package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/typed/apps/v1"
	"k8s.io/client-go/tools/clientcmd"
	"os"
)

func main() {

	kubeconfig := flag.String("kubeconfig", "/home/appscode/.kube/config", "location")

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)

	if err != nil {
		panic(err.Error())
	}
	// Create serializer for access api objects, here 'clientSet' is as serializer
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)

	fmt.Println("For Create Deployment")
	prompt()
	createDeploy(deploymentsClient)
	fmt.Println("For Update Deployment")
	prompt()
	updateDeploy(deploymentsClient)
	fmt.Println("For Delete Deployment")
	prompt()
	deleteDeploy(deploymentsClient)
}
func createDeploy(deploymentsClient v1.DeploymentInterface) {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "create-deployment",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(2),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "demo",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "demo",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "web",
							Image: "nginx:1.12",
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}
	fmt.Println("Creating deployment...")
	result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())
	fmt.Println()
}
func updateDeploy(deploymentsClient v1.DeploymentInterface) {
	deployName := "create-deployment"
	result, getErr := deploymentsClient.Get(context.TODO(), deployName, metav1.GetOptions{})
	if getErr != nil {
		panic(fmt.Errorf("Failed to get latest version of Deployment: %v", getErr))
	}
	result.Spec.Replicas = int32Ptr(3)                           // Increase replica count
	result.Spec.Template.Spec.Containers[0].Image = "nginx:1.13" // change nginx version
	_, updateErr := deploymentsClient.Update(context.TODO(), result, metav1.UpdateOptions{})

	if updateErr != nil {
		panic(fmt.Errorf("Update failed: %v", updateErr))
	}
	fmt.Println("Updated deployment...")
	fmt.Println()
}
func deleteDeploy(deploymentsClient v1.DeploymentInterface) {
	deployName := "create-deployment"
	fmt.Println("Deleting deployment...")
	deletePolicy := metav1.DeletePropagationForeground
	if err := deploymentsClient.Delete(context.TODO(), deployName, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		panic(err)
	}
	fmt.Println("Deleted deployment.")
	fmt.Println()
}
func prompt() {
	fmt.Printf("-> Press Return key to continue.")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		break
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	fmt.Println()
}
func int32Ptr(i int32) *int32 { return &i }
