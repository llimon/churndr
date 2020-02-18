package controller

import (
	"fmt"
	"os"
	"time"

	"github.com/llimon/churndr/common"
	clientset "github.com/llimon/churndr/pkg/generated/clientset/versioned"
	informers "github.com/llimon/churndr/pkg/generated/informers/externalversions"
	"github.com/llimon/churndr/pkg/signals"
	"github.com/llimon/churndr/util"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
)

func RunChurnNotifierController() error {
	var config *rest.Config
	var err error
	// set up signals so we handle the first shutdown signal gracefully
	stopCh := signals.SetupSignalHandler()

	if common.Config.InClusterConfiguration {
		config, err = rest.InClusterConfig()
		if err != nil {
			return err

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
			return err
		}

	}

	// creates the clientset
	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	// creates the clientset
	exampleClient, err := clientset.NewForConfig(config)
	if err != nil {
		return err
	}

	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClient, time.Second*30)
	exampleInformerFactory := informers.NewSharedInformerFactory(exampleClient, time.Second*30)

	controller := NewController(kubeClient, exampleClient,
		kubeInformerFactory.Apps().V1().Deployments(),
		exampleInformerFactory.Churndrcontroller().V1alpha1().Podchurns())
	// notice that there is no need to run Start methods in a separate goroutine. (i.e. go kubeInformerFactory.Start(stopCh)
	// Start method is non-blocking and runs all registered informers in a dedicated goroutine.
	kubeInformerFactory.Start(stopCh)
	exampleInformerFactory.Start(stopCh)

	if err = controller.Run(2, stopCh); err != nil {
		klog.Fatalf("Error running controller: %s", err.Error())
	}
	fmt.Println(stopCh, kubeInformerFactory, exampleInformerFactory)
	return nil
}
