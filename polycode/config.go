package polycode

import (
	"context"
	"encoding/json"
)

type AppConfig map[string]interface{}

func FromAppConfig(ctx context.Context, configObj any) {
	srvCtx, ok := ctx.(ServiceContext)
	if ok {
		ret := srvCtx.AppConfig()
		b, err := json.Marshal(configObj)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(b, &ret)
		if err != nil {
			panic(err)
		}
		return
	}
	wkfCtx, ok := ctx.(WorkflowContext)
	if ok {
		ret := wkfCtx.AppConfig()
		b, err := json.Marshal(configObj)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(b, &ret)
		if err != nil {
			panic(err)
		}
		return
	}
	panic("invalid context")
}
