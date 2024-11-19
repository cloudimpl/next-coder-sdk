package polycode

import (
	"context"
	"time"
)

type BaseContext interface {
	context.Context
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
	Function(function func(input any) (any, error)) Function
}

type ApiContext interface {
	WorkflowContext
}

type RawContext interface {
	ServiceContext
	WorkflowContext
	ServiceExec(req ExecServiceExtendedRequest) (ExecServiceResponse, error)
	ApiExec(req ExecApiExtendedRequest) (ExecApiResponse, error)
	DbGet(req QueryExtendedRequest) (map[string]interface{}, error)
	DbQuery(req QueryExtendedRequest) ([]map[string]interface{}, error)
	FileGet(req GetFileExtendedRequest) (GetFileResponse, error)
}

type ContextImpl struct {
	ctx           context.Context
	sessionId     string
	dataStore     DataStore
	fileStore     FileStore
	config        AppConfig
	serviceClient *ServiceClient
	logger        Logger
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

func (s ContextImpl) Function(function func(input any) (any, error)) Function {
	return Function{ctx: s.ctx, sessionId: s.sessionId, function: function, serviceClient: s.serviceClient}
}

func (s ContextImpl) ServiceExec(req ExecServiceExtendedRequest) (ExecServiceResponse, error) {
	return s.serviceClient.ExecServiceExtended(s.sessionId, req)
}

func (s ContextImpl) ApiExec(req ExecApiExtendedRequest) (ExecApiResponse, error) {
	return s.serviceClient.ExecApiExtended(s.sessionId, req)
}

func (s ContextImpl) DbGet(req QueryExtendedRequest) (map[string]interface{}, error) {
	return s.serviceClient.GetItemExtended(s.sessionId, req)
}

func (s ContextImpl) DbQuery(req QueryExtendedRequest) ([]map[string]interface{}, error) {
	return s.serviceClient.QueryItemsExtended(s.sessionId, req)
}

func (s ContextImpl) FileGet(req GetFileExtendedRequest) (GetFileResponse, error) {
	return s.serviceClient.GetFileExtended(s.sessionId, req)
}

func (s ContextImpl) Logger() Logger {
	return s.logger
}
