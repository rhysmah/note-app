// Package validator provides generic validation functionality for any type.
package validator

// ValidationRule is a function that performs a single validation check on type T.
// It returns an error if the validation fails, nil otherwise.
type ValidationRule[T any] func(*T) error

// Validator holds a collection of validation rules for type T.
type Validator[T any] struct {
	Rules []ValidationRule[T]
}

// NewValidator creates a new Validator instance for type T with an empty rule set.
func NewValidator[T any]() *Validator[T] {
	return &Validator[T]{
		Rules: []ValidationRule[T]{},
	}
}

// Run executes all validation rules in sequence.
// It returns the first error encountered or nil if all validations pass.
func (v *Validator[T]) Run(opts *T) error {
	for _, rule := range v.Rules {
		if err := rule(opts); err != nil {
			return err
		}
	}
	return nil
}