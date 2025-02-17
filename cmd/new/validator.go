package new

import (
	"fmt"
	"strings"

	"github.com/rhysmah/note-app/cmd/root"
	"github.com/rhysmah/note-app/validator"
)

func NewValidator() *validator.Validator[NewOptions] {
	return &validator.Validator[NewOptions] {
		Rules: []validator.ValidationRule[NewOptions] {
			validateNoteName,
		},
	}
}

func validateNoteName(opts *NewOptions) error {
	root.AppLogger.Start(fmt.Sprintf("Validating note name: '%s'", opts.noteName))

	noteNameTrimmed := strings.TrimSpace(opts.noteName)

	if len(noteNameTrimmed) > noteNameCharLimit {
		errMsg := fmt.Sprintf("name exceeds %d character limit", noteNameCharLimit)
		root.AppLogger.Fail(errMsg)
		return fmt.Errorf(errMsg)
	}

	if err := checkForIllegalCharacters(noteNameTrimmed); err != nil {
		return fmt.Errorf("invalid characters in note name: %w", err)
	}

	root.AppLogger.Success("Note name passed all validation checks")
	return nil
}

func checkForIllegalCharacters(noteName string) error {
	var illegalCharsFound []rune

	for _, char := range noteName {
		if strings.ContainsRune(illegalChars, char) {
			illegalCharsFound = append(illegalCharsFound, char)
		}
	}

	if len(illegalCharsFound) > 0 {
		errMsg := fmt.Sprintf("name contains illegal characters: %q", string(illegalCharsFound))
		root.AppLogger.Fail(errMsg)
		return fmt.Errorf(errMsg)
	}

	return nil
}