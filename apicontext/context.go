package apicontext

import (
	"context"
	"github.com/cloudimpl/next-coder-sdk/polycode"
)

func FromContext(ctx context.Context) (polycode.ApiContext, error) {
	value := ctx.Value("polycode.context")
	if value == nil {
		return polycode.ApiContext{}, polycode.ErrContextNotFound
	}

	return value.(polycode.ApiContext), nil
}
