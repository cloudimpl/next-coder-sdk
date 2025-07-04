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

type RemoteServiceBuilder struct {
	ctx           context.Context
	sessionId     string
	envId         string
	service       string
	serviceClient *ServiceClient
	tenantId      string
	partitionKey  string
}

func (r *RemoteServiceBuilder) WithTenantId(tenantId string) *RemoteServiceBuilder {
	r.tenantId = tenantId
	return r
}

func (r *RemoteServiceBuilder) WithPartitionKey(partitionKey string) *RemoteServiceBuilder {
	r.partitionKey = partitionKey
	return r
}

func (r *RemoteServiceBuilder) Get() RemoteService {
	return RemoteService{
		ctx:           r.ctx,
		sessionId:     r.sessionId,
		envId:         r.envId,
		service:       r.service,
		serviceClient: r.serviceClient,
		tenantId:      r.tenantId,
		partitionKey:  r.partitionKey,
	}
}

type RemoteService struct {
	ctx           context.Context
	sessionId     string
	envId         string
	service       string
	serviceClient *ServiceClient
	tenantId      string
	partitionKey  string
}

func (r RemoteService) RequestReply(options TaskOptions, method string, input any) Response {
	req := ExecServiceRequest{
		EnvId:        r.envId,
		Service:      r.service,
		TenantId:     r.tenantId,
		PartitionKey: r.partitionKey,
		Method:       method,
		Options:      options,
		Input:        input,
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
	req := ExecServiceRequest{
		EnvId:         r.envId,
		Service:       r.service,
		TenantId:      r.tenantId,
		PartitionKey:  r.partitionKey,
		Method:        method,
		Options:       options,
		FireAndForget: true,
		Input:         input,
	}

	output, err := r.serviceClient.ExecService(r.sessionId, req)
	if err != nil {
		fmt.Printf("client: exec task error: %v\n", err)
		return ErrTaskExecError.Wrap(err)
	}

	fmt.Printf("client: exec task output: %v\n", output)
	if output.IsError {
		return output.Error
	} else {
		return nil
	}
}

type RemoteAgentBuilder struct {
	ctx           context.Context
	sessionId     string
	envId         string
	agent         string
	serviceClient *ServiceClient
	tenantId      string
}

func (r *RemoteAgentBuilder) WithTenantId(tenantId string) *RemoteAgentBuilder {
	r.tenantId = tenantId
	return r
}

func (r *RemoteAgentBuilder) Get() RemoteAgent {
	return RemoteAgent{
		ctx:           r.ctx,
		sessionId:     r.sessionId,
		envId:         r.envId,
		agent:         r.agent,
		serviceClient: r.serviceClient,
		tenantId:      r.tenantId,
	}
}

type RemoteAgent struct {
	ctx           context.Context
	sessionId     string
	envId         string
	agent         string
	serviceClient *ServiceClient
	tenantId      string
}

func (r RemoteAgent) Call(options TaskOptions, input AgentInput) Response {
	req := ExecServiceRequest{
		EnvId:        r.envId,
		Service:      "agent-service",
		TenantId:     r.tenantId,
		PartitionKey: r.agent + ":" + input.SessionKey,
		Method:       "CallAgent",
		Options:      options,
		Headers: map[string]string{
			AgentNameHeader: r.agent,
		},
		Input: input,
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

type RemoteApp struct {
	ctx           context.Context
	sessionId     string
	envId         string
	appName       string
	serviceClient *ServiceClient
}

func (r RemoteApp) RequestReply(options TaskOptions, method string, input any) Response {
	req := ExecAppRequest{
		EnvId:   r.envId,
		AppName: r.appName,
		Method:  method,
		Options: options,
		Input:   input,
	}

	output, err := r.serviceClient.ExecApp(r.sessionId, req)
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

func (r RemoteApp) Send(options TaskOptions, method string, input any) error {
	req := ExecAppRequest{
		EnvId:         r.envId,
		AppName:       r.appName,
		Method:        method,
		Options:       options,
		FireAndForget: true,
		Input:         input,
	}

	output, err := r.serviceClient.ExecApp(r.sessionId, req)
	if err != nil {
		fmt.Printf("client: exec task error: %v\n", err)
		return ErrTaskExecError.Wrap(err)
	}

	fmt.Printf("client: exec task output: %v\n", output)
	if output.IsError {
		return output.Error
	} else {
		return nil
	}
}

type RemoteController struct {
	ctx           context.Context
	sessionId     string
	envId         string
	controller    string
	serviceClient *ServiceClient
}

func (r RemoteController) RequestReply(options TaskOptions, path string, apiReq ApiRequest) (ApiResponse, error) {
	req := ExecApiRequest{
		EnvId:      r.envId,
		Controller: r.controller,
		Path:       path,
		Options:    options,
		Request:    apiReq,
	}

	output, err := r.serviceClient.ExecApi(r.sessionId, req)
	if err != nil {
		return ApiResponse{}, err
	}

	if output.IsError {
		return ApiResponse{}, output.Error
	}

	return output.Response, nil
}

func (r RemoteController) Send(options TaskOptions, path string, apiReq ApiRequest) error {
	req := ExecApiRequest{
		EnvId:         r.envId,
		Controller:    r.controller,
		Path:          path,
		Options:       options,
		FireAndForget: true,
		Request:       apiReq,
	}

	output, err := r.serviceClient.ExecApi(r.sessionId, req)
	if err != nil {
		return err
	}

	if output.IsError {
		return output.Error
	}

	return nil
}

type Memo struct {
	ctx           context.Context
	sessionId     string
	serviceClient *ServiceClient
	getter        func() (any, error)
}

func (f Memo) Get() Response {
	req1 := ExecFuncRequest{
		Input: nil,
	}

	res1, err := f.serviceClient.ExecFunc(f.sessionId, req1)
	if err != nil {
		fmt.Printf("client: exec func error: %v\n", err)
		return Response{
			output:  nil,
			isError: true,
			error:   ErrTaskExecError.Wrap(err),
		}
	}

	if res1.IsCompleted {
		return Response{
			output:  res1.Output,
			isError: res1.IsError,
			error:   res1.Error,
		}
	}

	output, err := f.getter()
	var response Response
	if err != nil {
		response = Response{
			output:  nil,
			isError: true,
			error:   ErrTaskExecError.Wrap(err),
		}
	} else {
		response = Response{
			output:  output,
			isError: false,
			error:   Error{},
		}
	}

	req2 := ExecFuncResult{
		Input:   nil,
		Output:  response.output,
		IsError: response.isError,
		Error:   response.error,
	}

	err = f.serviceClient.ExecFuncResult(f.sessionId, req2)
	if err != nil {
		fmt.Printf("client: exec func result error: %v\n", err)
		return Response{
			output:  nil,
			isError: true,
			error:   ErrTaskExecError.Wrap(err),
		}
	}

	return response
}
