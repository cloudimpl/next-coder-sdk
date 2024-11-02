package polycode

import (
	"context"
	"encoding/json"
	"fmt"
)

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
	if output.IsError {
		return ThrowError(output.Error)
	}
	return FutureFrom(output.Output)
}

func (r RemoteService) Send(options TaskOptions, method string, input any) error {
	//TODO implement me
	panic("implement me")
}
