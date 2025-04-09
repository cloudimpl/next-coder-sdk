package polycode

import (
	"context"
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
	FileStore() FileStore
}

type WorkflowContext interface {
	BaseContext
	Service(service string) *RemoteServiceBuilder
	Controller(controller string) RemoteController
	Memo(getter func() (any, error)) Response
}

type ApiContext interface {
	WorkflowContext
}

type RawContext interface {
	BaseContext
	ServiceExec(req ExecServiceExtendedRequest) (ExecServiceResponse, error)
	ApiExec(req ExecApiExtendedRequest) (ExecApiResponse, error)
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

func (s ContextImpl) Controller(controller string) RemoteController {
	return RemoteController{ctx: s.ctx, sessionId: s.sessionId, controller: controller, serviceClient: s.serviceClient}
}

func (s ContextImpl) Memo(getter func() (any, error)) Response {
	m := Memo{ctx: s.ctx, sessionId: s.sessionId, getter: getter, serviceClient: s.serviceClient}
	return m.Get()
}

func (s ContextImpl) ServiceExec(req ExecServiceExtendedRequest) (ExecServiceResponse, error) {
	return s.serviceClient.ExecServiceExtended(s.sessionId, req)
}

func (s ContextImpl) ApiExec(req ExecApiExtendedRequest) (ExecApiResponse, error) {
	return s.serviceClient.ExecApiExtended(s.sessionId, req)
}

func (s ContextImpl) Logger() Logger {
	return s.logger
}

func (s ContextImpl) Acknowledge() error {
	return s.serviceClient.Acknowledge(s.sessionId)
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
