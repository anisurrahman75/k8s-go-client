package main

import (
	"context"
	"flag"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"time"
)

func main() {

	kubeconfig := flag.String("kubeconfig", "/home/appscode/.kube/config", "location")

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	//fmt.Println(config)
	if err != nil {
		panic(err.Error())
	}
	// Create serializer for access api objects, here 'clientSet' is as serializer
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	show_namespace_pods("default", *clientset)
	//show_namespace_deploy("default", *clientset)

}

func show_namespace_pods(namespace string, clientset kubernetes.Clientset) {

	for {
		pods, _ := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
		fmt.Printf("There are %d pods in %s cluster\n", len(pods.Items), namespace)
		for _, pod := range pods.Items {
			fmt.Println(pod.Name)
		}
		fmt.Println()
		time.Sleep(2 * time.Second)
	}
}
func show_namespace_deploy(namespace string, clientset kubernetes.Clientset) {

	for {
		dep, _ := clientset.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
		fmt.Printf("There are %d deploy in %s cluster\n", len(dep.Items), namespace)
		for _, d := range dep.Items {
			fmt.Println(d.Name)
		}
		fmt.Println()
		time.Sleep(2 * time.Second)
	}
}
