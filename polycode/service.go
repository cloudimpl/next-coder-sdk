package polycode

import (
	"context"
	"fmt"
)

type Response struct {
	output  any
	isError bool
	error   Error
}

func (r Response) IsError() bool {
	return r.isError
}

func (r Response) HasResult() bool {
	return r.output != nil
}

func (r Response) Get(ret any) error {
	if r.isError {
		return r.error
	}

	return ConvertType(r.output, ret)
}

func (r Response) GetAny() (any, error) {
	if r.isError {
		return nil, r.error
	} else {
		return r.output, nil
	}
}

type RemoteService struct {
	ctx           context.Context
	sessionId     string
	service       string
	serviceClient *ServiceClient
}

func (r RemoteService) RequestReply(options TaskOptions, method string, input any) Response {
	taskInput := TaskInput{
		Input: input,
	}

	req := ExecServiceRequest{
		Service: r.service,
		Method:  method,
		Options: options,
		Input:   taskInput,
	}

	output, err := r.serviceClient.ExecService(r.sessionId, req)
	if err != nil {
		fmt.Printf("client: exec task error: %v\n", err)
		return Response{
			output:  nil,
			isError: true,
			error:   ErrTaskExecError.Wrap(err),
		}
	}

	fmt.Printf("client: exec task output: %v\n", output)
	return Response{
		output:  output.Output,
		isError: output.IsError,
		error:   output.Error,
	}
}

func (r RemoteService) Send(options TaskOptions, method string, input any) error {
	panic("implement me")
}
