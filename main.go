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

	srv := schedular.NewServiceControl(conn)
	srv.ServiceExist("docker.elenet.me/sre/valid-python-sample:1")
}
