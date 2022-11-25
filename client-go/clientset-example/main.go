package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func echo(a any) {
	fmt.Println(a)
}

func echof(s string, a any) {
	fmt.Printf(s, a)
}

// 使用clientset 只能用内建资源

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

	echo("List nodes")
	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if errors.IsNotFound(err) {
		fmt.Printf("nodes not found %s\n", err)
	}
	for i := range nodes.Items {
		//echof("%+v", nodes.Items[i])
		echof("node: %s\n", nodes.Items[i].Name)
	}

	echo("\nList Pods\n")
	pods, err := clientset.CoreV1().Pods("develop").List(context.TODO(), metav1.ListOptions{})
	if errors.IsNotFound(err) {
		fmt.Printf("pods not found %s\n", err)
	}
	for _, v := range pods.Items {
		echof("pod: %s\n", v.Name)
	}

}
