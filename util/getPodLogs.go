package util

import (
	"bytes"
	"io"
	"os"

	"github.com/llimon/churndr/common"
	"github.com/llimon/churndr/common/util"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"k8s.io/client-go/rest"
)

func GetPreviousPodLogs(pod *corev1.Pod, containerName string, tailLines int64, limitBytes int64) (string, error) {
	var config *rest.Config
	var err error

	//tailLines := int64(200)
	//limitBytes := int64(4096)
	podLogOpts := corev1.PodLogOptions{
		Previous:   true,
		Container:  containerName,
		TailLines:  &tailLines,
		LimitBytes: &limitBytes,
	}
	if common.Config.InClusterConfiguration {
		config, err = rest.InClusterConfig()
		if err != nil {
			return "", err

		}
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			common.Sugar.Infof("Could not determine user home dir, setting it to /tmp")
			home = "/tmp"
		}
		kubeconfig := util.GetEnv("KUBECONFIG", home+"/.kube/config")
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
