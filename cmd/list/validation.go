package list

import (
	"fmt"

	"github.com/rhysmah/note-app/validator"
)

// NewValidator creates a validator with a predefined set of validation rules.
func NewValidator() *validator.Validator[ListOptions] {
	return &validator.Validator[ListOptions]{
		Rules: []validator.ValidationRule[ListOptions]{
			validateSortFieldExists,
			validateDateSortOrder,
			validateNameSortOrder,
			validateSortField,
			validateOrderField,
		},
	}
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
