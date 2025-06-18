package polycode

import (
	"errors"
	"fmt"
	"time"
)

type ParamWrapper struct {
	collection Collection
	exist      bool
	config     Config
}

func (p *ParamWrapper) IsExist() bool {
	return p.exist
}

func (p *ParamWrapper) Get() (string, error) {
	if !p.exist {
		return "", errors.New("config not found")
	}

	return p.config.Value, nil
}

func (p *ParamWrapper) Set(value string) error {
	p.config = Config{
		Id:        p.config.Id,
		Name:      p.config.Name,
		Value:     value,
		Version:   p.config.Version + 1,
		IsSecret:  p.config.IsSecret,
		Type:      p.config.Type,
		Scope:     p.config.Scope,
		Group:     p.config.Group,
		CreatedBy: p.config.CreatedBy,
		CreatedAt: p.config.CreatedAt,
		UpdatedAt: time.Now(),
	}

	err := p.collection.UpsertOne(p.config)
	if err != nil {
		return err
	}

	p.exist = true
	return nil
}

type ParamStore struct {
	collection Collection
}

func (p ParamStore) GlobalVar(group string, name string) (ParamWrapper, error) {
	config := Config{}
	exist, err := p.collection.GetOne(name, &config)
	if err != nil {
		return ParamWrapper{}, err
	}

	if exist {
		if config.Group != group {
			return ParamWrapper{}, errors.New("config group mismatch")
		}

		if config.IsSecret {
			return ParamWrapper{}, errors.New("config is secret")
		}

		return ParamWrapper{
			collection: p.collection,
			exist:      true,
			config:     config,
		}, nil
	} else {
		return ParamWrapper{
			collection: p.collection,
			exist:      false,
			config: Config{
				Id:        name,
				Name:      name,
				Value:     "",
				Version:   0,
				IsSecret:  false,
				Type:      "string",
				Scope:     Global,
				App:       "",
				Group:     group,
				CreatedBy: "",
				CreatedAt: time.Now(),
				UpdatedBy: "",
				UpdatedAt: time.Now(),
			},
		}, nil
	}
}

func (p ParamStore) GlobalSecret(group string, name string) (ParamWrapper, error) {
	config := Config{}
	exist, err := p.collection.GetOne(name, &config)
	if err != nil {
		return ParamWrapper{}, err
	}

	if exist {
		if config.Group != group {
			return ParamWrapper{}, errors.New("config group mismatch")
		}

		if !config.IsSecret {
			return ParamWrapper{}, errors.New("config is var")
		}

		return ParamWrapper{
			collection: p.collection,
			exist:      true,
			config:     config,
		}, nil
	} else {
		return ParamWrapper{
			collection: p.collection,
			exist:      false,
			config: Config{
				Id:        name,
				Name:      name,
				Value:     "",
				Version:   0,
				IsSecret:  true,
				Type:      "string",
				Scope:     Global,
				App:       "",
				Group:     group,
				CreatedBy: "",
				CreatedAt: time.Now(),
				UpdatedBy: "",
				UpdatedAt: time.Now(),
			},
		}, nil
	}
}

func (p ParamStore) AppVar(group string, name string) (ParamWrapper, error) {
	id := fmt.Sprintf("%s::%s", group, name)

	config := Config{}
	exist, err := p.collection.GetOne(id, &config)
	if err != nil {
		return ParamWrapper{}, err
	}

	if exist {
		if config.Group != group {
			return ParamWrapper{}, errors.New("config group mismatch")
		}

		if config.IsSecret {
			return ParamWrapper{}, errors.New("config is secret")
		}

		return ParamWrapper{
			collection: p.collection,
			exist:      true,
			config:     config,
		}, nil
	} else {
		return ParamWrapper{
			collection: p.collection,
			exist:      false,
			config: Config{
				Id:        id,
				Name:      name,
				Value:     "",
				Version:   0,
				IsSecret:  false,
				Type:      "string",
				Scope:     App,
				App:       GetClientEnv().AppName,
				Group:     group,
				CreatedBy: "",
				CreatedAt: time.Now(),
				UpdatedBy: "",
				UpdatedAt: time.Now(),
			},
		}, nil
	}
}

func (p ParamStore) AppSecret(group string, name string) (ParamWrapper, error) {
	id := fmt.Sprintf("%s::%s", group, name)

	config := Config{}
	exist, err := p.collection.GetOne(id, &config)
	if err != nil {
		return ParamWrapper{}, err
	}

	if exist {
		if config.Group != group {
			return ParamWrapper{}, errors.New("config group mismatch")
		}

		if !config.IsSecret {
			return ParamWrapper{}, errors.New("config is var")
		}

		return ParamWrapper{
			collection: p.collection,
			exist:      true,
			config:     config,
		}, nil
	} else {
		return ParamWrapper{
			collection: p.collection,
			exist:      false,
			config: Config{
				Id:        id,
				Name:      name,
				Value:     "",
				Version:   0,
				IsSecret:  false,
				Type:      "string",
				Scope:     App,
				App:       GetClientEnv().AppName,
				Group:     group,
				CreatedBy: "",
				CreatedAt: time.Now(),
				UpdatedBy: "",
				UpdatedAt: time.Now(),
			},
		}, nil
	}
}
