package polycode

import (
	"encoding/json"
)

type Future interface {
	Get(ret any) error
	GetAny() (any, error)
	IsNull() bool
}

func ThrowError(err error) Future {
	return ErrorFuture{error: err}
}

type ErrorFuture struct {
	error error
}

func (e ErrorFuture) GetAny() (any, error) {
	return nil, e.error
}

func (e ErrorFuture) IsNull() bool {
	return false
}

func (e ErrorFuture) Get(ret any) error {
	return e.error
}

type FutureImpl struct {
	data any
}

func (f FutureImpl) GetAny() (any, error) {
	return f.data, nil
}

func (f FutureImpl) IsNull() bool {
	return false
}

func (f FutureImpl) Get(ret any) error {
	switch ret.(type) {
	case map[string]interface{}:
		{
			ret = f.data
		}
	default:
		{
			b, err := json.Marshal(f.data)
			if err != nil {
				return err
			}
			switch ret.(type) {
			case *string:
				{
					*ret.(*string) = string(b)
					return nil
				}
			}
			return json.Unmarshal(b, ret)
		}
	}
	return nil
}

func FutureFrom(data any) Future {
	return FutureImpl{data: data}
}

var nullFuture Future = NullFutureImpl{}

func NullFuture() Future {
	return nullFuture
}

type NullFutureImpl struct {
}

func (n NullFutureImpl) GetAny() (any, error) {
	return nil, nil
}

func (n NullFutureImpl) IsNull() bool {
	return true
}

func (n NullFutureImpl) Get(ret any) error {
	return nil
}
