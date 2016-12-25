package schedular

import (
	"errors"

	"strings"

	"github.com/lijianying10/log"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

type ServiceControl struct {
	Conn         *Connection
	K8SNameSpace string
}

func NewServiceControl(conn *Connection, NameSpace string) *ServiceControl {
	var serviceControl ServiceControl
	serviceControl.Conn = conn
	serviceControl.K8SNameSpace = NameSpace
	return &serviceControl
}

func (sc *ServiceControl) ServiceNameHandle(imageName string) (string, string, error) {
	strs := strings.Split(imageName, ":")
	if len(strs) != 2 {
		return "", "", errors.New("error define image name")
	}

	strs[0] = strings.Replace(strs[0], "/", "-", -1)
	//strs[0] = strings.Replace(strs[0], "-", "_", -1)
	strs[0] = strings.Replace(strs[0], ".", "-", -1)

	return strs[0], strs[1], nil
}

func (sc *ServiceControl) ServiceExist(imageName string) (bool, error) {
	var err error
	img, ver, err := sc.ServiceNameHandle(imageName)
	if err != nil {
		return false, err
	}
	pods, err := sc.Conn.ClientSet.Core().Pods(sc.K8SNameSpace).List(v1.ListOptions{
		LabelSelector: "servicetype=ewf,image=" + img + ",version=" + ver,
	})
	if err != nil {
		log.Error("get images error", err.Error())
		return false, errors.New("error get service")
	}

	if len(pods.Items) == 0 {
		return false, nil
	}

	if len(pods.Items) >= 1 {
		return false, errors.New("Multi service error")
	}

	if pods.Items[0].Status.Phase != "Running" {
		return false, errors.New("Service is not running")
	}

	return true, nil
}

func (sc *ServiceControl) ServiceCreate(imageName string) error {
	img, ver, err := sc.ServiceNameHandle(imageName)
	if err != nil {
		return err
	}
	_, err = sc.Conn.ClientSet.ExtensionsV1beta1Client.Deployments(sc.K8SNameSpace).Create(&v1beta1.Deployment{
		ObjectMeta: v1.ObjectMeta{
			Name: img + "-" + ver + "ewf",
			Labels: map[string]string{
				"servicetype": "ewf",
				"image":       img,
				"version":     ver,
			},
			Namespace: sc.K8SNameSpace,
		},
		Spec: v1beta1.DeploymentSpec{
			Replicas: 2,
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						v1.Container{
							Name:  img + "-" + ver + "-ewfrunc",
							Image: imageName,
						},
					},
				},
			},
		},
	})
	if err != nil {
		return err
	}

	_, err = sc.Conn.ClientSet.CoreV1Client.Services(sc.K8SNameSpace).Create(&v1.Service{
		ObjectMeta: v1.ObjectMeta{
			Name: img + "-" + ver + "ewf",
			Labels: map[string]string{
				"servicetype": "ewf",
				"image":       img,
				"version":     ver,
			},
			Namespace: sc.K8SNameSpace,
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				v1.ServicePort{
					Protocol:   v1.ProtocolTCP,
					Port:       80,
					TargetPort: 10080,
				},
			},
			Selector: map[string]string{
				"servicetype": "ewf",
				"image":       img,
				"version":     ver,
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}
