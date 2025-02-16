package list

import "fmt"

// ValidationRule defines a function that validates a specific aspect of ListOptions.
type ValidationRule func(*ListOptions) error

// validator contains a collection of validation rules to be executed.
type validator struct {
	rules []ValidationRule
}

// NewValidator creates a validator with a predefined set of validation rules.
func NewValidator() *validator {
	return &validator{
		rules: []ValidationRule{
			validateSortFieldExists,
			validateDateSortOrder,
			validateNameSortOrder,
			validateSortField,
			validateOrderField,
		},
	}
}

// Run executes all validation rules in sequence, returning the first error encountered.
func (v *validator) Run(opts *ListOptions) error {
	for _, rule := range v.rules {
		if err := rule(opts); err != nil {
			return err
		}
	}
	return nil
}

// validateSortFieldExists checks if a sort field has been specified.
func validateSortFieldExists(opts *ListOptions) error {
	if opts.SortField == "" {
		return fmt.Errorf("%q (%q) flag required with `list`. Available sort fields: %q",
			sortByCmd, sortByCmdShort, availableSortFields())
	}
	return nil
}

// validateNameSortOrder ensures name sorting uses alphabetical or reverse alphabetical order.
func validateNameSortOrder(opts *ListOptions) error {
	if opts.SortField == SortFieldName {
		if opts.SortOrder != SortOrderAlph && opts.SortOrder != SortOrderRAlph {
			return fmt.Errorf("when sorting by name, order must be either %q or %q, got %q",
				SortOrderAlph, SortOrderRAlph, opts.SortOrder)
		}
	}
	return nil
}

// validateDateSortOrder ensures date-based sorting uses newest or oldest order.
func validateDateSortOrder(opts *ListOptions) error {
	if opts.SortField == SortFieldCreated || opts.SortField == SortFieldModified {
		if opts.SortOrder != SortOrderNewest && opts.SortOrder != SortOrderOldest {
			return fmt.Errorf("when sorting by date, order must be either %q or %q, got %q",
				SortOrderNewest, SortOrderOldest, opts.SortOrder)
		}
	}
	return nil
}

// validateSortField verifies the sort field is one of the predefined valid options.
func validateSortField(opts *ListOptions) error {
	if _, valid := sortFieldDescriptions[opts.SortField]; !valid {
		return fmt.Errorf("invalid sort field: %q.\nValid sort fields: %q", opts.SortField, availableSortFields())
	}
	return nil
}

// validateOrderField verifies the sort order is one of the predefined valid options.
func validateOrderField(opts *ListOptions) error {
	if _, valid := sortOrderDescriptions[opts.SortOrder]; !valid {
		return fmt.Errorf("invalid sort order selected: %q. Valid sort orders: %q", opts.SortOrder, availableSortOrders())
	}
	return nil
}
