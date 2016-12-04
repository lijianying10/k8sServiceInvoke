package main

import (
	"flag"
	"fmt"

	"github.com/lijianying10/k8sServiceInvoke/schedular"
)

var (
	kubeconfig = flag.String("kubeconfig", "/home/a/.kube/config", "absolute path to the kubeconfig file")
)

func main() {
	fmt.Println("starting")
	conn := schedular.NewConnection(*kubeconfig)

	K8SNameSpace := "sre"

	srv := schedular.NewServiceControl(conn, K8SNameSpace)
	exist, err := srv.ServiceExist("docker.elenet.me/sre/valid-python-sample:1")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("exist: ", exist)
	err = srv.ServiceCreate("docker.elenet.me/sre/valid-python-sample:1")
	if err != nil {
		fmt.Println(err.Error())
	}
}
