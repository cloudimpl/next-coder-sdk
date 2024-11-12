package polycode

var currentValidator Validator = DummyValidator{}

func SetValidator(v Validator) {
	currentValidator = v
}

type Validator interface {
	Validate(obj any) error
}

type DummyValidator struct {
}

func (v DummyValidator) Validate(obj any) error {
	return nil
}
