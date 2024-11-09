package rawcontext

import (
	"context"
	"github.com/cloudimpl/next-coder-sdk/polycode"
)

func FromContext(ctx context.Context) (polycode.RawContext, error) {
	ctxImpl, ok := ctx.(polycode.ContextImpl)
	if ok {
		return ctxImpl, nil
	}

	value := ctx.Value("polycode.context")
	if value == nil {
		return nil, polycode.ErrContextNotFound
	}

	return value.(polycode.ContextImpl), nil
}
