package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/CloudImpl-Inc/next-coder-sdk/client/db"
	"github.com/CloudImpl-Inc/next-coder-sdk/polycode"
	"os"
	"time"
)

type RuntimeSupport interface {
	GetRuntime() polycode.Runtime
	InvokeWorkflow(workflowContext polycode.WorkflowContext, input polycode.TaskInput) (any, error)
}

type ServiceContext struct {
	ctx       context.Context
	sessionId string
	db        polycode.Database
	fileStore polycode.FileStore
	option    polycode.TaskOptions
}

func (s ServiceContext) Option() polycode.TaskOptions {
	return s.option
}

func (s ServiceContext) Deadline() (deadline time.Time, ok bool) {
	return s.ctx.Deadline()
}

func (s ServiceContext) Done() <-chan struct{} {
	return s.ctx.Done()
}

func (s ServiceContext) Err() error {
	return s.ctx.Err()
}

func (s ServiceContext) Value(key any) any {
	return s.ctx.Value(key)
}

func (s ServiceContext) Db() polycode.Database {
	return s.db
}

func (s ServiceContext) FileStore() polycode.FileStore {
	return s.fileStore
}

type WorkflowContext struct {
	ctx           context.Context
	sessionId     string
	serviceClient *ServiceClient
}

func (A WorkflowContext) Deadline() (deadline time.Time, ok bool) {
	return A.ctx.Deadline()
}

func (A WorkflowContext) Done() <-chan struct{} {
	return A.ctx.Done()
}

func (A WorkflowContext) Err() error {
	return A.ctx.Err()
}

func (A WorkflowContext) Value(key any) any {
	return A.ctx.Value(key)
}

func (A WorkflowContext) Service(serviceId string) (polycode.RemoteService, error) {

	return remoteService{ctx: A.ctx, sessionId: A.sessionId, serviceId: serviceId, serviceClient: A.serviceClient}, nil
}

func (A WorkflowContext) LocalService() (polycode.RemoteService, error) {
	return remoteService{ctx: A.ctx, sessionId: A.sessionId, serviceId: os.Getenv("polycode_SERVICE_ID"), serviceClient: A.serviceClient}, nil
}

type remoteService struct {
	ctx           context.Context
	sessionId     string
	serviceId     string
	serviceClient *ServiceClient
}

func (r remoteService) RequestReply(options polycode.TaskOptions, method string, input any) polycode.Future {

	remoteCtx := RemoteTaskContext{
		Context:    r.ctx,
		SessionId:  r.sessionId,
		ServiceId:  r.serviceId,
		EntryPoint: method,
		Options:    options,
	}

	b, err := json.Marshal(input)
	if err != nil {
		return polycode.ThrowError(err)
	}
	taskInput := polycode.TaskInput{
		NoArg:     false,
		TargetReq: string(b),
	}

	output, err := r.serviceClient.ExecTask(remoteCtx, method, taskInput)
	if err != nil {
		println(fmt.Sprintf("execTask error %s", err.Error()))
		return polycode.ThrowError(err)
	}
	println(fmt.Sprintf("exec task output %v", output))
	if output.Error != nil {
		return polycode.ThrowError(output.Error)
	}
	return polycode.FutureFrom(output.Output)
}

func (r remoteService) Send(options polycode.TaskOptions, method string, input any) error {
	//TODO implement me
	panic("implement me")
}

func NewRuntimeContext(ctx context.Context, db db.Database) RuntimeContext {
	return RuntimeContext{
		ctx: ctx,
		db:  db,
	}
}
