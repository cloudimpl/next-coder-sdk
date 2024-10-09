package polycode

import (
	"context"
)

type ServiceContext interface {
	context.Context
	Db() Database
	AppConfig() AppConfig
	FileStore() FileStore
	Option() TaskOptions
}

type ReadOnlyServiceContext interface {
	context.Context
	Db() ReadOnlyDatabase
	FileStore() ReadOnlyFileStore
	Option() TaskOptions
}

type WorkflowContext interface {
	context.Context
	Service(serviceId string) (RemoteService, error)
	AppConfig() AppConfig
	LocalService() (RemoteService, error)
}

func OptionNone() TaskOptions {
	return TaskOptions{}
}

func ToWorkflowContext(ctx context.Context) (WorkflowContext, error) {
	ret, ok := ctx.Value("polycode.context").(WorkflowContext)
	if !ok {
		return nil, ErrInvalidContext
	}
	return ret, nil
}
