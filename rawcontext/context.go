package apicontext

import (
	"context"
	"github.com/cloudimpl/next-coder-sdk/polycode"
)

func FromContext(ctx context.Context) (polycode.RawContext, error) {
	rawContext, ok := ctx.(polycode.RawContext)
	if !ok {
		return nil, polycode.ErrContextNotFound
	}

	return rawContext, nil
}
