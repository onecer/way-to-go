package main

import (
	"flag"
	"fmt"
	"path/filepath"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	stopCh := make(chan struct{})
	defer close(stopCh)

	shardInformers := informers.NewSharedInformerFactory(clientset, time.Minute)
	informer := shardInformers.Core().V1().Pods().Informer()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			mObj := obj.(metav1.Object)
			fmt.Println("Pod added to store: ", mObj.GetName())
			fmt.Println("add")
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			mObj := newObj.(metav1.Object)
			fmt.Println("Pod updated in store: ", mObj.GetName())
			fmt.Println("update")
		},
		DeleteFunc: func(obj interface{}) {
			mObj := obj.(metav1.Object)
			fmt.Println("Pod deleted from store: ", mObj.GetName())
			fmt.Println("delete")
		},
	})
	go informer.Run(stopCh)
}
