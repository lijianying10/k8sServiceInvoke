package schedular

import (
	"errors"

	"strings"

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

func (sc *ServiceControl) ServiceNameHandle(imageName string) (string, string, error) {
	strs := strings.Split(imageName, ":")
	if len(strs) != 2 {
		return "", "", errors.New("error define image name")
	}

	strs[0] = strings.Replace(strs[0], "/", "_", -1)

	return strs[0], strs[1], nil
}

func (sc *ServiceControl) ServiceExist(imageName string) (bool, error) {
	var err error
	img, ver, err := sc.ServiceNameHandle(imageName)
	if err != nil {
		return false, err
	}
	pods, err := sc.Conn.ClientSet.Core().Pods("").List(v1.ListOptions{
		LabelSelector: "servicetype=ewf,image=" + img + ",version=" + ver,
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
