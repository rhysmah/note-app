package validator

type ValidationRule[T any] func(*T) error

type Validator[T any] struct {
	rules []ValidationRule[T]
}

func NewValidator[T any](rules ...ValidationRule[T]) *Validator[T] {
	return &Validator[T]{
		rules: rules,
	}
}

func (v *Validator[T]) Run(opts *T) error {
	for _, rule := range v.rules {
		if err := rule(opts); err != nil {
			return err
		}
	}
	return nil
}