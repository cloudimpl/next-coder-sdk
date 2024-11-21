package polycode

import (
	"context"
	"encoding/json"
	"errors"
)

type AppConfig map[string]interface{}

func FromAppConfig(ctx context.Context, configObj any) error {
	baseCtx, ok := ctx.(BaseContext)
	if ok {
		ret := baseCtx.AppConfig()
		b, err := json.Marshal(configObj)
		if err != nil {
			return err
		}

		err = json.Unmarshal(b, &ret)
		if err != nil {
			return err
		}
		return nil
	}

	return errors.New("invalid context")
}
