package main

import (
    "flag"
    "fmt"
    "os"
    "path/filepath"
    "time"
	"context"

    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
)

func main() {
    var kubeconfig *string
	var ctx, _ = context.WithCancel(context.Background())
    if home := homeDir(); home != "" {
        kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
    } else {
        kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
    }
    flag.Parse()

    // uses the current context in kubeconfig
    config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
    if err != nil {
        panic(err.Error())
    }

    // creates the clientset
    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        panic(err.Error())
    }
		
	namespaceList, err := clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	for _, namespace := range namespaceList.Items {
//		deploymentList, err := clientset.CoreV1().Deployments.List(metav1.ListOptions{})

		deploymentsClient := clientset.AppsV1().Deployments(namespace.Name)
		deploymentList, err := deploymentsClient.List(ctx, metav1.ListOptions{})
		if err != nil {
			panic(err)
		}

		for _, deploymentName := range deploymentList.Items {
			clientset := clientset.AppsV1().Deployments(namespace.Name)
			deployment, err := clientset.Get(ctx, deploymentName.Name, metav1.GetOptions{})
			if err != nil {
				panic(nil)
			}	
			name := deployment.GetName()
			fmt.Println("------Deployment Name------")
			fmt.Println("name ->", name)
			fmt.Println("")
			containers := &deployment.Spec.Template.Spec.Containers
			for i := range *containers {
				c := *containers
				fmt.Println("------Container Name && Image Name------")
				fmt.Println(c[i].Name, "  ", c[i].Image)
				fmt.Println("\n")
			}
				
				
//	        pods, err := clientset.CoreV1().Pods(namespace.Name).List(ctx, metav1.ListOptions{})
//	        if err != nil {
//	            panic(err.Error())
//	        }
//	        fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
		}	
		time.Sleep(10 * time.Second)
	}
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
    }
    return os.Getenv("USERPROFILE") // windows
}
