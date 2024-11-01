package polycode

import (
	"context"
	"time"
)

type ServiceContext struct {
	ctx       context.Context
	sessionId string
	dataStore DataStore
	fileStore FileStore
	option    TaskOptions
	config    AppConfig
}

func (s ServiceContext) AppConfig() AppConfig {
	return s.config
}

func (s ServiceContext) Option() TaskOptions {
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

func (s ServiceContext) Db() DataStore {
	return s.dataStore
}

func (s ServiceContext) FileStore() FileStore {
	return s.fileStore
}

type WorkflowContext struct {
	ctx           context.Context
	sessionId     string
	serviceClient *ServiceClient
	config        AppConfig
}

func (wc WorkflowContext) AppConfig() AppConfig {
	return wc.config
}

func (wc WorkflowContext) Deadline() (deadline time.Time, ok bool) {
	return wc.ctx.Deadline()
}

func (wc WorkflowContext) Done() <-chan struct{} {
	return wc.ctx.Done()
}

func (wc WorkflowContext) Err() error {
	return wc.ctx.Err()
}

func (wc WorkflowContext) Value(key any) any {
	return wc.ctx.Value(key)
}

func (wc WorkflowContext) Service(serviceId string) (RemoteService, error) {
	return RemoteService{ctx: wc.ctx, sessionId: wc.sessionId, serviceId: serviceId, serviceClient: wc.serviceClient}, nil
}

type ApiContext struct {
	ctx           context.Context
	sessionId     string
	serviceClient *ServiceClient
	config        AppConfig
}

func (wc ApiContext) AppConfig() AppConfig {
	return wc.config
}

func (wc ApiContext) Deadline() (deadline time.Time, ok bool) {
	return wc.ctx.Deadline()
}

func (wc ApiContext) Done() <-chan struct{} {
	return wc.ctx.Done()
}

func (wc ApiContext) Err() error {
	return wc.ctx.Err()
}

func (wc ApiContext) Value(key any) any {
	return wc.ctx.Value(key)
}

func (wc ApiContext) Service(serviceId string) (RemoteService, error) {
	return RemoteService{ctx: wc.ctx, sessionId: wc.sessionId, serviceId: serviceId, serviceClient: wc.serviceClient}, nil
}
