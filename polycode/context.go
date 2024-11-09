package polycode

import (
	"context"
	"time"
)

type ServiceContext interface {
	AppConfig() AppConfig
	Db() DataStore
	FileStore() FileStore
}

type WorkflowContext interface {
	AppConfig() AppConfig
	Service(service string) (RemoteService, error)
	Controller(controller string) (RemoteController, error)
}

type ApiContext interface {
	AppConfig() AppConfig
	Service(service string) (RemoteService, error)
	Controller(controller string) (RemoteController, error)
}

type RawContext interface {
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

func (s ContextImpl) Service(service string) (RemoteService, error) {
	return RemoteService{ctx: s.ctx, sessionId: s.sessionId, service: service, serviceClient: s.serviceClient}, nil
}

func (s ContextImpl) Controller(controller string) (RemoteController, error) {
	return RemoteController{ctx: s.ctx, sessionId: s.sessionId, controller: controller, serviceClient: s.serviceClient}, nil
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
