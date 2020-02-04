package controller

import (
	"fmt"

	"github.com/kubernetes/client-go/informers"

	"github.com/llimon/churndr/common"
	"github.com/llimon/churndr/util"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"k8s.io/client-go/tools/cache"
)

const (
	K8S_APP = "k8s-app"
	APP     = "app"
)

// This should be done by reading a environment variable or a file
var excludeApps = []string{"kiamon", "kiam", "moot"}

/*KubeGetPods get details for  kubernetes resources names matching
the filter*/
func KubeGetPods() (map[string]interface{}, error) {
	resources := map[string]interface{}{}
	// creates the in-cluster config
	/*
		currconfig, err := rest.InClusterConfig()
		if err != nil {
			fmt.Println(err.Error())
			return resources, err
		}
	*/
	// use the current context in kubeconfig
	//kubeconfig := "/Users/llimon/.kube/config"
	kubeconfig := util.GetEnv("KUBECONFIG", "/Users/llimon/.kube/config")
	currconfig, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(currconfig)
	if err != nil {
		common.Sugar.Infof("Error getting clientset of kubernetes")
		return resources, err
	}
	// Create the shared informer factory and use the client to connect to
	// Kubernetes
	factory := informers.NewSharedInformerFactory(clientset, 0)

	// Get the informer for the right resource, in this case a Pod
	informer := factory.Core().V1().Pods().Informer()

	// Create a channel to stops the shared informer gracefully
	stopper := make(chan struct{})
	defer close(stopper)

	// Kubernetes serves an utility to handle API crashes
	defer runtime.HandleCrash()

	// This is the part where your custom code gets triggered based on the
	// event that the shared informer catches
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		// When a new pod gets created
		AddFunc: onAdd,
		// When a pod gets updated
		UpdateFunc: onUpdate,
		// When a pod gets deleted
		DeleteFunc: onDelete,
	})

	// You need to start the informer, in my case, it runs in the background
	go informer.Run(stopper)

	if !cache.WaitForCacheSync(stopper, informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return resources, nil
	}
	<-stopper

	return resources, nil
}

func podHasFilteredLabel(pod *corev1.Pod, label string, excludes []string) (string, error) {
	//var selector labels.Selector

	if len(pod.Labels) == 0 {
		return "", fmt.Errorf("pod %v because has no labels", pod.Name)
	} else {
		val, ok := pod.Labels[APP]
		if ok {
			if util.Contains(excludes, val) {
			} else {
				return val, nil
			}
		}
	}

	return "", fmt.Errorf("Pod %v does not have label %v", pod.Name, APP)

}

func onAdd(obj interface{}) {
	// Cast the obj as node
	pod := obj.(*corev1.Pod)

	label, ok := podHasFilteredLabel(pod, APP, excludeApps)
	if ok == nil {
		common.Sugar.Infow("onAdd",
			"pod name", pod.ObjectMeta.Name,
			"Label", label,
		)
	}

	// need to fill in database on start of controller.

}

func onUpdate(obj interface{}, obj2 interface{}) {
	// Cast the obj as node
	pod := obj.(*corev1.Pod)
	/*
		    _, ok := pod.GetLabels()[]
		    if ok {
		        fmt.Printf("It has the label!")
			}
	*/

	/*
		    &ContainerStateRunning{StartedAt:2020-01-29 14:24:03 -0800 PST,} nil
			&ContainerStateTerminated{ExitCode:0,Signal:0,Reason:Completed,Message:,StartedAt:2020-01-29 14:24:03 -0800 PST,FinishedAt:2020-01-29 14:24:17 -0800 PST,ContainerID:docker://f80ee5fd323a44f0a86dc9ea4101535409faafa63298bea637f55cf835ff17a2,}
	*/

	_, ok := podHasFilteredLabel(pod, APP, []string{})
	if ok == nil {
		var containerList []common.ContainerDB

		for _, containerStatus := range pod.Status.ContainerStatuses {
			if containerStatus.State.Running != nil {

				common.Sugar.Infow("OnUpdate",
					"Namespace", pod.ObjectMeta.Namespace,
					"POD", pod.ObjectMeta.Name,
					"Container", containerStatus.Name,
					"Running", util.GetDateString(containerStatus.State.Running.StartedAt.Time),
					"Restarts", containerStatus.RestartCount,
				)

				containerList = append(containerList, common.ContainerDB{
					Running: &corev1.ContainerStateRunning{
						StartedAt: containerStatus.State.Running.StartedAt,
					},
					Name:         containerStatus.Name,
					RestartCount: containerStatus.RestartCount,
				})

			} else if containerStatus.State.Terminated != nil {
				common.Sugar.Infow("OnUpdate",
					"Namespace", pod.ObjectMeta.Namespace,
					"POD", pod.ObjectMeta.Name,
					"Container", containerStatus.Name,
					"FinishedAt", util.GetDateString(containerStatus.State.Terminated.FinishedAt.Time),
					"ExitCode", containerStatus.State.Terminated.ExitCode,
					"Signal", containerStatus.State.Terminated.Signal,
					"Reason", containerStatus.State.Terminated.Reason,
					"Restarts", containerStatus.RestartCount,
				)
				containerList = append(containerList, common.ContainerDB{
					Terminated: &corev1.ContainerStateTerminated{
						FinishedAt: containerStatus.State.Terminated.FinishedAt,
						ExitCode:   containerStatus.State.Terminated.ExitCode,
						Signal:     containerStatus.State.Terminated.Signal,
						Reason:     containerStatus.State.Terminated.Reason,
					},
					Name:         containerStatus.Name,
					RestartCount: containerStatus.RestartCount,
				})

				if _, found := common.PodCache[pod.ObjectMeta.Name]; found == false {
					fmt.Println("Marking pod deletion on our pod database cache")
				}
			} else if containerStatus.State.Waiting != nil {
				common.Sugar.Infow("OnUpdate",
					"Namespace", pod.ObjectMeta.Namespace,
					"POD", pod.ObjectMeta.Name,
					"Container", containerStatus.Name,
				)
			}
			//newPod.container = containerList
			//fmt.Println(containerStatus.Name, containerStatus.RestartCount, containerStatus.State.Running, containerStatus.State.Terminated)
		}
		// Persist lastTimeReported
		ltr := common.PodCache[pod.ObjectMeta.Name].LastTimeReported
		common.PodCache[pod.ObjectMeta.Name] = common.PodDB{
			Name:             pod.ObjectMeta.Name,
			Namespace:        pod.ObjectMeta.Namespace,
			Container:        containerList,
			LastTimeReported: ltr,
		}
		fmt.Println(containerList)

	}

}

func onDelete(obj interface{}) {
	//fmt.Println("onDelete")
	// Cast the obj as node
	pod := obj.(*corev1.Pod)
	label, ok := podHasFilteredLabel(pod, APP, excludeApps)
	if ok == nil {
		fmt.Printf("onDelete-POD %v has label %v=%v", pod.ObjectMeta.Name, APP, label)
		if _, found := common.PodCache[pod.ObjectMeta.Name]; found == false {
			fmt.Println("Marking pod deletion on our pod database cache")
			// BUG/FIX: Some how mark it as deleted
			common.PodCache[pod.ObjectMeta.Name] = common.PodDB{
				Name:      pod.ObjectMeta.Name,
				Namespace: pod.ObjectMeta.Namespace,
			}
		}
	}

}
