package hybrid

import (
	"github.com/CloudImpl-Inc/next-coder-sdk/polycode"
	"github.com/CloudImpl-Inc/next-coder-sdk/runtime/aws"
	"github.com/CloudImpl-Inc/next-coder-sdk/runtime/local"
	"os"
)

type Runtime struct {
	runtime polycode.Runtime
}

func (r Runtime) AppConfig() polycode.AppConfig {
	return r.runtime.AppConfig()
}

func (r Runtime) Name() string {
	return "Hybrid"
}

func (r Runtime) Start(params []any) error {
	if os.Getenv("AWS_LAMBDA_RUNTIME_API") != "" {
		r.runtime = &aws.Runtime{}
	} else {
		r.runtime = &local.Runtime{}
	}
	return r.runtime.Start(params)
}
