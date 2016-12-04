package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/lijianying10/k8sServiceInvoke/schedular"
	"k8s.io/client-go/pkg/api/v1"
)

var (
	kubeconfig = flag.String("kubeconfig", "/home/a/.kube/config", "absolute path to the kubeconfig file")
)

func main() {
	fmt.Println("starting")
	conn := schedular.NewConnection(*kubeconfig)

	for {
		pods, err := conn.ClientSet.Core().Pods("").List(v1.ListOptions{
			LabelSelector: "app=pcmysql,pod-template-hash=1667360149",
		})
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
		time.Sleep(10 * time.Second)
	}
}
