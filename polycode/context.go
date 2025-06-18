package polycode

import (
	"context"
	"time"
)

type AuthContext struct {
	claims map[string]interface{}
}

type BaseContext interface {
	context.Context
	Meta() ContextMeta
	AppConfig() AppConfig
	AuthContext() AuthContext
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
	Signal(signalName string) Signal
	RealtimeChannel(channelName string) RealtimeChannel
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

func (s ContextImpl) Signal(signalName string) Signal {
	return Signal{
		name:          signalName,
		sessionId:     s.sessionId,
		serviceClient: s.serviceClient,
	}
}

func (s ContextImpl) RealtimeChannel(channelName string) RealtimeChannel {
	return RealtimeChannel{
		name:          channelName,
		sessionId:     s.sessionId,
		serviceClient: s.serviceClient,
	}
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
