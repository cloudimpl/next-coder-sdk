package polycode

import (
	"context"
	"encoding/json"
	"errors"
	"time"
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

type ConfigScope uint64

const (
	Global ConfigScope = iota
	App
)

type Config struct {
	Id        string      `polycode:"id" json:"id"`
	Name      string      `json:"name"`
	Value     string      `json:"value"`
	Version   uint64      `json:"version"`
	IsSecret  bool        `json:"isSecret"`
	Type      string      `json:"type"`
	Scope     ConfigScope `json:"scope"`
	Group     string      `json:"group"`
	App       string      `json:"app"`
	CreatedBy string      `json:"createdBy"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedBy string      `json:"updatedBy"`
	UpdatedAt time.Time   `json:"updatedAt"`
}

type GetConfigRequest struct {
	Id string `json:"id"`
}

type ConfigGroup struct {
	Name string `polycode:"id" json:"name"`
}

type AuditTrail struct {
	Id         string    `polycode:"id" json:"id"`
	ConfigId   string    `json:"configId"`
	ConfigName string    `json:"configName"`
	Action     string    `json:"action"`
	OldValue   string    `json:"oldValue"`
	NewValue   string    `json:"newValue"`
	User       string    `json:"user"`
	Timestamp  time.Time `json:"timestamp"`
	IsSecret   bool      `json:"isSecret"`
}

type SaveConfigRequest struct {
	Name     string      `json:"name"`
	Value    string      `json:"value"`
	IsSecret bool        `json:"isSecret"`
	Type     string      `json:"type"`
	Scope    ConfigScope `json:"scope"`
	Group    string      `json:"group"`
	App      string      `json:"app"`
}

type ListConfigGroupsRequest struct {
}

type CreateConfigGroupRequest struct {
	Name string `polycode:"id" json:"name"`
}

type ListConfigRequest struct {
	Name  string      `json:"name"`
	Scope ConfigScope `json:"scope"`
	Group string      `json:"group"`
	App   string      `json:"app"`
}
