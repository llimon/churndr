package controller

import (
	"fmt"
	"os"
	"reflect"

	//"github.com/kubernetes/client-go/informers"
	"k8s.io/client-go/informers"

	"github.com/llimon/churndr/common"
	"github.com/llimon/churndr/util"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
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
	var currconfig *rest.Config
	var err error
	resources := map[string]interface{}{}

	// creates the in-cluster config
	if common.Config.InClusterConfiguration {
		currconfig, err = rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			common.Sugar.Infof("Could not determine user home dir, setting it to /tmp")
			home = "/tmp"
		}
		kubeconfig := util.GetEnv("KUBECONFIG", home+"/.kube/config")
		currconfig, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			panic(err.Error())
		}

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

// Return true if pod is in one of the namespaces we monitor
func podInMonitoredNamespace(pod *corev1.Pod) bool {

	for _, v := range common.Config.Namespaces {
		if pod.ObjectMeta.Namespace == v {
			return true
		}
	}
	return false
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

	//label, ok := podHasFilteredLabel(pod, APP, excludeApps)
	ok := podInMonitoredNamespace(pod)
	if ok {
		common.Sugar.Infow("onAdd",
			"pod name", pod.ObjectMeta.Name,
		)
	}

	// need to fill in database on start of controller.
}

func onUpdate(obj interface{}, obj2 interface{}) {
	doSavePodStatus := false
	// Cast the obj as pod
	pod := obj.(*corev1.Pod)

	ok := podInMonitoredNamespace(pod)
	if ok {

		// Only process Pod types we care about.
		isValidKind := false
		for _, o := range pod.OwnerReferences {
			if o.Kind == "ReplicaSet" {
				isValidKind = true
			}
		}
		if !isValidKind {
			common.Sugar.Infow("Skip pod",
				"Namespace", pod.Namespace,
				"Name", pod.Name)
			return
		}
		// Detect pods stuck in Pending for longer than N time
		if pod.Status.Phase == "Pending" { //&& pod.Status.StartTime.Time {

		}
		var containerList []common.ContainerDB
		for _, containerStatus := range pod.Status.ContainerStatuses {
			// if containerStatus.State.Running != nil {

			// 	common.Sugar.Infow("OnUpdate",
			// 		"Namespace", pod.ObjectMeta.Namespace,
			// 		"POD", pod.ObjectMeta.Name,
			// 		"Container", containerStatus.Name,
			// 		"Running", util.GetDateString(containerStatus.State.Running.StartedAt.Time),
			// 		"Restarts", containerStatus.RestartCount,
			// 	)

			// 	containerList = append(containerList, common.ContainerDB{
			// 		Running: &corev1.ContainerStateRunning{
			// 			StartedAt: containerStatus.State.Running.StartedAt,
			// 		},
			// 		Name:         containerStatus.Name,
			// 		RestartCount: containerStatus.RestartCount,
			// 	})
			// 	// BUG/FIX: We want to make sure if there is a restart that we verity StartedAt Time and see if we need to report or not
			// 	if containerStatus.RestartCount > 0 {
			// 		doSavePodStatus = true
			// 	}

			//} else if containerStatus.State.Terminated != nil {
			if containerStatus.State.Terminated != nil {
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
				c := common.ContainerDB{
					Terminated: &corev1.ContainerStateTerminated{
						FinishedAt: containerStatus.State.Terminated.FinishedAt,
						ExitCode:   containerStatus.State.Terminated.ExitCode,
						Signal:     containerStatus.State.Terminated.Signal,
						Reason:     containerStatus.State.Terminated.Reason,
					},
					Name:         containerStatus.Name,
					RestartCount: containerStatus.RestartCount,
					Logs:         "/pod/log/container/" + pod.ObjectMeta.Namespace + "/" + pod.ObjectMeta.Name + "/" + containerStatus.Name + "/" + fmt.Sprint(containerStatus.RestartCount),
				}
				c.SetLink("logs", "/pod/log/container/"+pod.ObjectMeta.Namespace+"/"+pod.ObjectMeta.Name+"/"+containerStatus.Name+"/"+fmt.Sprint(containerStatus.RestartCount), "")
				containerList = append(containerList, c)
				doSavePodStatus = true

				// Save logs of terminated pod for posterity
				if containerStatus.RestartCount > 0 {
					go util.PersistPodLogs(pod, containerStatus.Name, containerStatus.RestartCount)
				}

				// if _, found := common.PodCache[pod.ObjectMeta.Name]; found == false {
				// 	fmt.Println("Marking pod deletion on our pod database cache")
				// }
				//} else if containerStatus.State.Waiting != nil {
			} else {
				// BUG/FIX: We don't do anything for Waiting. There might be some opportunity to catch issues with pods on waiting stage.
				//          Need to review.
				common.Sugar.Infow("OnUpdate",
					"Namespace", pod.ObjectMeta.Namespace,
					"POD", pod.ObjectMeta.Name,
					"Container", containerStatus.Name,
				)
			}
			//newPod.container = containerList
			//fmt.Println(containerStatus.Name, containerStatus.RestartCount, containerStatus.State.Running, containerStatus.State.Terminated)
		}
		if doSavePodStatus {
			// Persist lastTimeReported
			ltr := common.PodCache[pod.ObjectMeta.Name].LastTimeReported
			p := common.PodDB{
				Name:             pod.ObjectMeta.Name,
				Namespace:        pod.ObjectMeta.Namespace,
				Container:        containerList,
				LastTimeReported: ltr,
			}
			p.SetLink("self", "/pod/"+pod.ObjectMeta.Namespace+"/"+pod.ObjectMeta.Name, "")
			common.PodCache[pod.ObjectMeta.Name] = p

		}
	}

}

func onDelete(obj interface{}) {
	//fmt.Println("onDelete")
	objReflect := reflect.ValueOf(obj)
	// BUG/FIX: somehow obj == cache.DeletedFinalStateUnknown  happens. Should we monitor for this condition?

	// BUG/FIX: How do I detect pods terminated because of eviction?
	fmt.Println(">>>>>>> objReflect.Type().String()", objReflect.Type().String())
	if objReflect.Type().String() == "*v1.Pod" {
		// cast object as pod
		pod := obj.(*corev1.Pod)
		ok := podInMonitoredNamespace(pod)
		if ok {
			fmt.Printf("onDelete-POD %v has label %v", pod.ObjectMeta.Name, APP)
			if _, found := common.PodCache[pod.ObjectMeta.Name]; found == false {
				fmt.Println("Marking pod deletion on our pod database cache")
				//BUG/FIX; this is incomplete, for now don't monitor for it.
				// common.PodCache[pod.ObjectMeta.Name] = common.PodDB{
				// 	Name:      pod.ObjectMeta.Name,
				// 	Namespace: pod.ObjectMeta.Namespace,
				// }
			}
		}
	}

}
