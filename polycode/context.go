package polycode

import (
	"context"
	"time"
)

type AuthContext struct {
	Claims map[string]interface{} `json:"claims"`
}

type BaseContext interface {
	context.Context
	Meta() ContextMeta
	AppConfig() AppConfig
	AuthContext() AuthContext
	Logger() Logger
	ParamStore() ParamStore
	UnsafeDb() *UnsafeDataStoreBuilder
	FileStore() FileStore
}

type ServiceContext interface {
	BaseContext
	Db() DataStore
}

type WorkflowContext interface {
	BaseContext
	Service(service string) *RemoteServiceBuilder
	Agent(agent string) *RemoteAgentBuilder
	ServiceEx(envId string, service string) *RemoteServiceBuilder
	AgentEx(envId string, agent string) *RemoteAgentBuilder
	App(appName string) RemoteApp
	AppEx(envId string, appName string) RemoteApp
	Controller(controller string) RemoteController
	ControllerEx(envId string, controller string) RemoteController
	Memo(getter func() (any, error)) Response
	Signal(signalName string) Signal
	ClientChannel(channelName string) ClientChannel
	Lock(key string) Lock
}

type ApiContext interface {
	WorkflowContext
}

type RawContext interface {
	BaseContext
	GetMeta(group string, typeName string, key string) (map[string]interface{}, error)
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
	authCtx       AuthContext
}

func (s ContextImpl) Meta() ContextMeta {
	return s.meta
}

func (s ContextImpl) AuthContext() AuthContext {
	return s.authCtx
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

func (s ContextImpl) ParamStore() ParamStore {
	return ParamStore{
		collection: s.Db().Collection("config"),
	}
}

func (s ContextImpl) FileStore() FileStore {
	return s.fileStore
}

func (s ContextImpl) Service(service string) *RemoteServiceBuilder {
	return &RemoteServiceBuilder{
		ctx: s.ctx, sessionId: s.sessionId, service: service, serviceClient: s.serviceClient,
	}
}

func (s ContextImpl) Agent(agent string) *RemoteAgentBuilder {
	return &RemoteAgentBuilder{
		ctx: s.ctx, sessionId: s.sessionId, agent: agent, serviceClient: s.serviceClient,
	}
}

func (s ContextImpl) ServiceEx(envId string, service string) *RemoteServiceBuilder {
	return &RemoteServiceBuilder{
		ctx: s.ctx, sessionId: s.sessionId, envId: envId, service: service, serviceClient: s.serviceClient,
	}
}

func (s ContextImpl) AgentEx(envId string, agent string) *RemoteAgentBuilder {
	return &RemoteAgentBuilder{
		ctx: s.ctx, sessionId: s.sessionId, envId: envId, agent: agent, serviceClient: s.serviceClient,
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

func (s ContextImpl) Signal(signalName string) Signal {
	return Signal{
		name:          signalName,
		sessionId:     s.sessionId,
		serviceClient: s.serviceClient,
	}
}

func (s ContextImpl) ClientChannel(channelName string) ClientChannel {
	return ClientChannel{
		name:          channelName,
		sessionId:     s.sessionId,
		serviceClient: s.serviceClient,
	}
}

func (s ContextImpl) Lock(key string) Lock {
	return Lock{
		client:    s.serviceClient,
		sessionId: s.sessionId,
		key:       key,
	}
}

func (s ContextImpl) GetMeta(group string, typeName string, key string) (map[string]interface{}, error) {
	req := GetMetaDataRequest{
		Group: group,
		Type:  typeName,
		Key:   key,
	}

	return s.serviceClient.GetMeta(s.sessionId, req)
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
