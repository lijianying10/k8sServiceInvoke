package main

import (
	"flag"
	"fmt"

	"os/user"

	"github.com/eleme/esm-agent/log"
	"github.com/lijianying10/k8sServiceInvoke/schedular"
)

var (
	kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
)

func main() {
	fmt.Println("starting")
	if *kubeconfig == "" {
		user, err := user.Current()
		if err != nil {
			log.Fatal("error get current client user")
		}
		*kubeconfig = user.HomeDir + "/.kube/config"
	}
	conn := schedular.NewConnection(*kubeconfig)

	K8SNameSpace := "sre"

	srv := schedular.NewServiceControl(conn, K8SNameSpace)
	exist, err := srv.ServiceExist("docker.elenet.me/sre/golang-ewf-web-service:1")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("exist: ", exist)
	err = srv.ServiceCreate("docker.elenet.me/sre/golang-ewf-web-service:1")
	if err != nil {
		fmt.Println(err.Error())
	}
}
