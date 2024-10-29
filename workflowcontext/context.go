package workflowcontext

import (
	"context"
	"github.com/cloudimpl/next-coder-sdk/polycode"
)

func FromContext(ctx context.Context) (polycode.WorkflowContext, error) {
	value := ctx.Value("polycode.context")
	if value == nil {
		return polycode.WorkflowContext{}, polycode.ErrContextNotFound
	}

	return value.(polycode.WorkflowContext), nil
}
