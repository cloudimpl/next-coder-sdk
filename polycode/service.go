package polycode

import (
	"fmt"
)

var ErrServiceNotFound = DefineError("polycode.service", 1, "service not found,serviceId: [%s]")

type NoArg string
type ServiceRegistry struct {
	services map[string]Service
}

var serviceRegistry = ServiceRegistry{services: make(map[string]Service)}

func GetRemoteService(serviceId string) (Service, error) {
	service, ok := serviceRegistry.services[serviceId]
	if !ok {
		return nil, ErrServiceNotFound.With(serviceId)
	}
	return service, nil
}

var mainService Service = nil

func RegisterService(service Service) {
	fmt.Println("register service ", service.GetName())
	if mainService != nil {
		panic("main service already set")
	}
	mainService = service
}

func GetService() (Service, error) {
	if mainService == nil {
		return nil, fmt.Errorf("service not set")
	}
	return mainService, nil
}

type Service interface {
	GetName() string
	GetInputType(method string) (any, error)
	ExecuteService(ctx ServiceContext, method string, input any) (any, error)
	ExecuteWorkflow(ctx WorkflowContext, method string, input any) (any, error)
	IsWorkflow(method string) bool
}
