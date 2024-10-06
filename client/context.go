package client

import (
	"context"
	"github.com/CloudImpl-Inc/next-coder-sdk/client/db"
	"github.com/CloudImpl-Inc/next-coder-sdk/polycode"
	"time"
)

type RemoteTaskContext struct {
	Context    context.Context      `json:"-"`
	SessionId  string               `json:"sessionId"`
	ServiceId  string               `json:"serviceId"`
	EntryPoint string               `json:"entryPoint"`
	Options    polycode.TaskOptions `json:"options"`
}

type RuntimeContext struct {
	ctx context.Context
	db  db.Database
}

func (r RuntimeContext) Deadline() (deadline time.Time, ok bool) {
	return r.ctx.Deadline()
}

func (r RuntimeContext) Done() <-chan struct{} {
	return r.ctx.Done()
}

func (r RuntimeContext) Err() error {
	return r.ctx.Err()
}

func (r RuntimeContext) Value(key any) any {
	return r.ctx.Value(key)
}

func (r RuntimeContext) GetDB() db.Database {
	return r.db
}

const (
	EventAPICall EventType = iota
	EventServiceCall
	InternalCall
)

type EventType int
type Event struct {
	Type      EventType
	Context   TaskContext
	TaskInput polycode.TaskInput
	Error     string
}
