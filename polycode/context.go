package polycode

import (
	"context"
	"fmt"
	"time"
)

type BaseContext interface {
	context.Context
	Meta() ContextMeta
	AppConfig() AppConfig
	Logger() Logger
}

type ServiceContext interface {
	BaseContext
	Db() DataStore
	UnsafeDb() *UnsafeDataStoreBuilder
	FileStore() FileStore
}

type WorkflowContext interface {
	BaseContext
	Service(service string) *RemoteServiceBuilder
	ServiceEx(envId string, service string) *RemoteServiceBuilder
	App(appName string) RemoteApp
	AppEx(envId string, appName string) RemoteApp
	Controller(controller string) RemoteController
	ControllerEx(envId string, controller string) RemoteController
	UnsafeDb() *UnsafeDataStoreBuilder
	Memo(getter func() (any, error)) Response
	SignalAwait(signalName string) Response
	SignalResumeSuccess(taskId string, signalName string, data any) error
	SignalResumeError(taskId string, signalName string, err Error) error
	EmitRealtimeEvent(channel string, data any) error
}

type ApiContext interface {
	WorkflowContext
}

type RawContext interface {
	BaseContext
	Counter(group string, name string, ttl int64) Counter
}

type ContextImpl struct {
	ctx           context.Context
	sessionId     string
	dataStore     DataStore
	fileStore     FileStore
	config        AppConfig
	serviceClient *ServiceClient
	logger        Logger
	meta          ContextMeta
}

func (s ContextImpl) Meta() ContextMeta {
	return s.meta
}

func (s ContextImpl) AppConfig() AppConfig {
	return s.config
}

func (s ContextImpl) Deadline() (deadline time.Time, ok bool) {
	return s.ctx.Deadline()
}

func (s ContextImpl) Done() <-chan struct{} {
	return s.ctx.Done()
}

func (s ContextImpl) Err() error {
	return s.ctx.Err()
}

func (s ContextImpl) Value(key any) any {
	return s.ctx.Value(key)
}

func (s ContextImpl) Db() DataStore {
	return s.dataStore
}

func (s ContextImpl) FileStore() FileStore {
	return s.fileStore
}

func (s ContextImpl) Service(service string) *RemoteServiceBuilder {
	return &RemoteServiceBuilder{
		ctx: s.ctx, sessionId: s.sessionId, service: service, serviceClient: s.serviceClient,
	}
}

func (s ContextImpl) ServiceEx(envId string, service string) *RemoteServiceBuilder {
	return &RemoteServiceBuilder{
		ctx: s.ctx, sessionId: s.sessionId, envId: envId, service: service, serviceClient: s.serviceClient,
	}
}

func (s ContextImpl) App(appName string) RemoteApp {
	return RemoteApp{
		ctx: s.ctx, sessionId: s.sessionId, appName: appName, serviceClient: s.serviceClient,
	}
}

func (s ContextImpl) AppEx(envId string, appName string) RemoteApp {
	return RemoteApp{
		ctx: s.ctx, sessionId: s.sessionId, envId: envId, appName: appName, serviceClient: s.serviceClient,
	}
}

func (s ContextImpl) Controller(controller string) RemoteController {
	return RemoteController{ctx: s.ctx, sessionId: s.sessionId, controller: controller, serviceClient: s.serviceClient}
}

func (s ContextImpl) ControllerEx(envId string, controller string) RemoteController {
	return RemoteController{ctx: s.ctx, sessionId: s.sessionId, envId: envId, controller: controller, serviceClient: s.serviceClient}
}

func (s ContextImpl) UnsafeDb() *UnsafeDataStoreBuilder {
	return &UnsafeDataStoreBuilder{
		client: s.serviceClient, sessionId: s.sessionId,
	}
}

func (s ContextImpl) Memo(getter func() (any, error)) Response {
	m := Memo{ctx: s.ctx, sessionId: s.sessionId, getter: getter, serviceClient: s.serviceClient}
	return m.Get()
}

func (s ContextImpl) Logger() Logger {
	return s.logger
}

func (s ContextImpl) Acknowledge() error {
	return s.serviceClient.Acknowledge(s.sessionId)
}

func (s ContextImpl) SignalAwait(signalName string) Response {
	req := SignalWaitRequest{
		SignalName: signalName,
	}

	output, err := s.serviceClient.WaitForSignal(s.sessionId, req)
	if err != nil {
		fmt.Printf("client: signal await error: %v\n", err)
		return Response{
			output:  nil,
			isError: true,
			error:   ErrTaskExecError.Wrap(err),
		}
	}

	fmt.Printf("client: signal await output: %v\n", output)
	return Response{
		output:  output.Output,
		isError: output.IsError,
		error:   output.Error,
	}
}

func (s ContextImpl) SignalResumeSuccess(taskId string, signalName string, data any) error {
	req := SignalEmitRequest{
		TaskId:     taskId,
		SignalName: signalName,
		Output:     data,
		IsError:    false,
	}

	return s.serviceClient.EmitSignal(s.sessionId, req)
}

func (s ContextImpl) SignalResumeError(taskId string, signalName string, err Error) error {
	req := SignalEmitRequest{
		TaskId:     taskId,
		SignalName: signalName,
		IsError:    true,
		Error:      err,
	}

	return s.serviceClient.EmitSignal(s.sessionId, req)
}

func (s ContextImpl) EmitRealtimeEvent(channel string, data any) error {
	req := RealtimeEventEmitRequest{
		Channel: channel,
		Input:   data,
	}

	return s.serviceClient.EmitRealtimeEvent(s.sessionId, req)
}

func (s ContextImpl) Counter(group string, name string, ttl int64) Counter {
	return Counter{
		client:    s.serviceClient,
		sessionId: s.sessionId,
		group:     group,
		name:      name,
		ttl:       ttl,
	}
}
