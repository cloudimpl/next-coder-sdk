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
		service:       r.service,
		serviceClient: r.serviceClient,
		tenantId:      r.tenantId,
		partitionKey:  r.partitionKey,
	}
}

type RemoteService struct {
	ctx           context.Context
	sessionId     string
	service       string
	serviceClient *ServiceClient
	tenantId      string
	partitionKey  string
}

func (r RemoteService) RequestReply(options TaskOptions, method string, input any) Response {
	req := ExecServiceRequest{
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
	panic("implement me")
}

type RemoteController struct {
	ctx           context.Context
	sessionId     string
	controller    string
	serviceClient *ServiceClient
}

func (r RemoteController) RequestReply(options TaskOptions, path string, apiReq ApiRequest) (ApiResponse, error) {
	req := ExecApiRequest{
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
