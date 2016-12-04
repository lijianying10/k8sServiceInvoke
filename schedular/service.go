package schedular

import (
	"errors"

	"github.com/lijianying10/log"
	"k8s.io/client-go/pkg/api/v1"
)

type ServiceControl struct {
	Conn *Connection
}

func NewServiceControl(conn *Connection) *ServiceControl {
	var serviceControl ServiceControl
	serviceControl.Conn = conn
	return &serviceControl
}

func (sc *ServiceControl) ServiceExist(imageName string) (bool, error) {
	pods, err := sc.Conn.ClientSet.Core().Pods("").List(v1.ListOptions{
		LabelSelector: imageName,
	})
	if err != nil {
		log.Error("get images error", err.Error())
		return false, errors.New("error get service")
	}

	if len(pods.Items) != 1 {
		return false, errors.New("Multi service error")
	}

	if pods.Items[0].Status.Phase != "Running" {
		return false, errors.New("Service is not running")
	}

	return true, nil
}
