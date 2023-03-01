package main

import (
	"context"
	"flag"
	"fmt"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	kubeconfig := flag.String("kubeconfig", "/home/appscode/.kube/config", "location")

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	watcher, err := clientset.CoreV1().Pods(apiv1.NamespaceDefault).Watch(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	for event := range watcher.ResultChan() {
		p := event.Object.(*apiv1.Pod)
		switch event.Type {
		case watch.Added:
			fmt.Printf("Pods %s/%s added\n", p.ObjectMeta.Namespace, p.ObjectMeta.Name)
		case watch.Modified:
			fmt.Printf("Pods %s/%s modified\n", p.ObjectMeta.Namespace, p.ObjectMeta.Name)
		case watch.Deleted:
			fmt.Printf("Pods %s/%s deleted\n", p.ObjectMeta.Namespace, p.ObjectMeta.Name)
		}
	}
}
