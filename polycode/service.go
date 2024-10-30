package polycode

import (
	"context"
	"encoding/json"
	"fmt"
)

var serviceMap = make(map[string]Service)

func RegisterService(service Service) {
	fmt.Println("register service ", service.GetName())
	serviceMap[service.GetName()] = service
}

func GetService(serviceName string) (Service, error) {
	service := serviceMap[serviceName]
	if service == nil {
		return nil, fmt.Errorf("service not set")
	}
	return service, nil
}

type Service interface {
	GetName() string
	GetInputType(method string) (any, error)
	ExecuteService(ctx ServiceContext, method string, input any) (any, error)
	ExecuteWorkflow(ctx WorkflowContext, method string, input any) (any, error)
	IsWorkflow(method string) bool
}

type RemoteService struct {
	ctx           context.Context
	sessionId     string
	serviceId     string
	serviceClient *ServiceClient
}

func (r RemoteService) RequestReply(options TaskOptions, method string, input any) Future {
	b, err := json.Marshal(input)
	if err != nil {
		return ThrowError(err)
	}

	taskInput := TaskInput{
		NoArg:     false,
		TargetReq: string(b),
	}
	req := ExecRequest{
		ServiceId:  r.serviceId,
		EntryPoint: method,
		Options:    options,
		Input:      taskInput,
	}

	output, err := r.serviceClient.ExecService(r.sessionId, req)
	if err != nil {
		println(fmt.Sprintf("execTask error %s", err.Error()))
		return ThrowError(err)
	}
	println(fmt.Sprintf("exec task output %v", output))
	if output.Error != nil {
		return ThrowError(output.Error)
	}
	return FutureFrom(output.Output)
}

func (r RemoteService) Send(options TaskOptions, method string, input any) error {
	//TODO implement me
	panic("implement me")
}
