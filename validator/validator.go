package validator

type ValidationRule[T any] func(*T) error

type Validator[T any] struct {
	Rules []ValidationRule[T]
}

func NewValidator[T any]() *Validator[T] {
	return &Validator[T]{
		Rules: []ValidationRule[T]{},
	}
}

func (v *Validator[T]) Run(opts *T) error {
	for _, rule := range v.Rules {
		if err := rule(opts); err != nil {
			return err
		}
	}
	return nil
}
