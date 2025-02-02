package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const (
	illegalChars      string = "\\/:*?\"<>|: ."
	noteNameCharLimit int    = 50
)

func init() {
	rootCmd.AddCommand(new)
}

var new = &cobra.Command{
	Use:   "create",
	Short: "Create a new note",
	Long: `Create a new note with the specified name.
The note will be saved as '[note-name]_[date].txt' in your notes directory.
Note names cannot contain special characters or exceed 50 characters.`,

	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(1)(cmd, args); err != nil {
			return fmt.Errorf("you must specify a name when creating a note")
		}
		return nil
	},

	Run: func(cmd *cobra.Command, args []string) {
		appLogger.Start(fmt.Sprintf("Creating new note with name: '%s'", args[0]))

		notesDir := dirManager.NotesDir()
		if notesDir == "" {
			appLogger.Fail("Cannot access notes directory")
			return
		}

		err := validateNoteName(args[0])
		if err != nil {
			fmt.Printf("Name validation failed for '%s': %v", args[0], err)
			return
		}
		appLogger.Success(fmt.Sprintf("Name validation passed for '%s'", args[0]))

		err = createAndSaveNote(notesDir, args[0])
		if err != nil {
			fmt.Printf("Note creation failed: %v", err)
			return
		}
		appLogger.End("Note creation process completed successfully")
	},
}

func createAndSaveNote(notesDirPath, noteName string) error {
	appLogger.Start(fmt.Sprintf("Creating note file for: '%s'", noteName))

	fullNoteName := noteName + "_" + time.Now().Format("2006_01_02_15_04") + ".txt"
	notePathFinal := filepath.Join(notesDirPath, fullNoteName)

	appLogger.Info(fmt.Sprintf("Target note path: %s", notePathFinal))

	if _, err := os.Stat(notePathFinal); err == nil {
		errMsg := fmt.Sprintf("Note already exists: %s", fullNoteName)
		appLogger.Fail(errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	file, err := os.Create(notePathFinal)
	if err != nil {
		errMsg := fmt.Sprintf("File creation failed: %v", err)
		appLogger.Fail(errMsg)
		return fmt.Errorf("problem creating note: %v", err)
	}
	defer file.Close()

	successMsg := fmt.Sprintf("Note created at: %s", notePathFinal)
	appLogger.Success(successMsg)
	fmt.Printf("Note successfully created at %s\n", notePathFinal)
	return nil
}

func validateNoteName(noteName string) error {
	appLogger.Start(fmt.Sprintf("Validating note name: '%s'", noteName))

	if len(noteName) > noteNameCharLimit {
		errMsg := fmt.Sprintf("Name exceeds %d character limit", noteNameCharLimit)
		appLogger.Fail(errMsg)
		return fmt.Errorf("note name cannot exceed %d characters", noteNameCharLimit)
	}

	if strings.TrimSpace(noteName) != noteName {
		errMsg := "Name contains leading or trailing spaces"
		appLogger.Fail(errMsg)
		return fmt.Errorf("note name cannot begin or end with spaces")
	}

	var illegalCharsFound []rune
	for _, char := range noteName {
		if strings.ContainsRune(illegalChars, char) {
			illegalCharsFound = append(illegalCharsFound, char)
		}
	}

	if len(illegalCharsFound) > 0 {
		errMsg := fmt.Sprintf("Name contains illegal characters: %q", string(illegalCharsFound))
		appLogger.Fail(errMsg)
		return fmt.Errorf("note name contains illegal characters: %q", string(illegalCharsFound))
	}

	appLogger.Success("Note name validation passed all checks")
	return nil
}
