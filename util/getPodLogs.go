package util

import (
	"bytes"
	"io"

	"github.com/llimon/churndr/common"
	"github.com/llimon/churndr/common/util"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"k8s.io/client-go/rest"
)

func getPodLogs(pod corev1.Pod) (string, error) {
	var config *rest.Config
	var err error
	podLogOpts := corev1.PodLogOptions{}
	if common.Config.InClusterConfiguration {
		config, err = rest.InClusterConfig()
		if err != nil {
			return "", err

		}
	} else {
		kubeconfig := util.GetEnv("KUBECONFIG", "/Users/llimon/.kube/config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return "", err
		}

	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return "", err
	}
	req := clientset.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &podLogOpts)
	podLogs, err := req.Stream()
	if err != nil {
		return "", err
	}
	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		return "", err
	}
	str := buf.String()

	return str, nil
}
