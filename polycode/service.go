package polycode

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/CloudImpl-Inc/next-coder-sdk/client"
)

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

type RemoteService struct {
	ctx           context.Context
	sessionId     string
	serviceId     string
	serviceClient *client.ServiceClient
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
	req := client.ExecRequest{
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
