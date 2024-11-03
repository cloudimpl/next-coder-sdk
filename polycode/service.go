package polycode

import (
	"context"
	"fmt"
)

type Response struct {
	output any
	error  any
}

func (r Response) IsError() bool {
	return r.error != nil
}

func (r Response) HasResult() bool {
	return r.output != nil
}

func (r Response) Get(ret any) error {
	if r.error != nil {
		return r.error.(error)
	}

	return ConvertType(r.output, ret)
}

func (r Response) GetAny() (any, error) {
	if r.error != nil {
		return nil, r.error.(error)
	} else {
		return r.output, nil
	}
}

type RemoteService struct {
	ctx           context.Context
	sessionId     string
	serviceId     string
	serviceClient *ServiceClient
}

func (r RemoteService) RequestReply(options TaskOptions, method string, input any) (Response, error) {
	taskInput := TaskInput{
		Input: input,
	}

	req := ExecRequest{
		ServiceId:  r.serviceId,
		EntryPoint: method,
		Options:    options,
		Input:      taskInput,
	}

	output, err := r.serviceClient.ExecService(r.sessionId, req)
	if err != nil {
		fmt.Printf("client: exec task error: %v\n", err)
		return Response{}, err
	}

	fmt.Printf("client: exec task output: %v\n", output)
	return Response{
		output: output.Output,
		error:  output.Error,
	}, nil
}

func (r RemoteService) Send(options TaskOptions, method string, input any) error {
	panic("implement me")
}
