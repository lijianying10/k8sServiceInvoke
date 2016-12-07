package schedular

import (
	"github.com/lijianying10/log"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Connection struct {
	Config    *rest.Config
	ClientSet *kubernetes.Clientset
}

func NewConnection(ConfigPath string) *Connection {
	var Conn Connection
	var err error
	Conn.Config, err = clientcmd.BuildConfigFromFlags("", ConfigPath)
	if err != nil {
		log.Fatal("error read config: ", err.Error())
	}

	Conn.ClientSet, err = kubernetes.NewForConfig(Conn.Config)
	if err != nil {
		log.Fatal("error read config: ", err.Error())
	}

	return &Conn
}
