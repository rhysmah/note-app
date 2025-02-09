package new

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rhysmah/note-app/cmd/root"
	"github.com/spf13/cobra"
)

const (
	illegalChars      string = "\\/:*?\"<>|: ."
	noteNameCharLimit int    = 50
)

func init() {
	root.RootCmd.AddCommand(new)
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
		root.AppLogger.Start(fmt.Sprintf("Creating new note with name: '%s'", args[0]))

		if err := createNote(args[0]); err != nil {
			fmt.Printf("There was an error creating your note: %v\n", err)
			os.Exit(1)
		}

		root.AppLogger.End("Note creation process completed successfully")
	},
}

func createNote(noteName string) error {
	notesDir := root.DirManager.NotesDir()
	if notesDir == "" {
		errMsg := fmt.Sprintf("Cannot access notes directory: %s", notesDir)
		root.AppLogger.Fail(errMsg)
		return fmt.Errorf(errMsg)
	}

	if err := validateNoteName(noteName); err != nil {
		return fmt.Errorf("invalid note name: %w", err)
	}

	if err := createAndSaveNote(noteName, notesDir); err != nil {
		return fmt.Errorf("failed to create note %s: %w", noteName, err)
	}

	return nil
}

func createAndSaveNote(noteName, notesDirPath string) error {
	root.AppLogger.Start(fmt.Sprintf("Creating note '%s' in directory %s...", noteName, notesDirPath))

	fullNoteName := noteName + "_" + time.Now().Format("2006_01_02_15_04") + ".txt"
	notePath := filepath.Join(notesDirPath, fullNoteName)

	// Check if note already exists
	if _, err := os.Stat(notePath); err == nil {
		errMsg := fmt.Sprintf("note %q already exists", fullNoteName)
		root.AppLogger.Fail(errMsg)
		return fmt.Errorf(errMsg)
	}

	// Create note
	file, err := os.Create(notePath)
	if err != nil {
		errMsg := fmt.Sprintf("failed to create file: %v", err)
		root.AppLogger.Fail(errMsg)
		return fmt.Errorf(errMsg)
	}
	defer file.Close()

	successMsg := fmt.Sprintf("note created at: %s", notePath)
	root.AppLogger.Success(successMsg)
	fmt.Printf("Created note: %s\n", notePath)
	return nil
}

func validateNoteName(noteName string) error {
	root.AppLogger.Start(fmt.Sprintf("Validating note name: '%s'", noteName))

	noteNameTrimmed := strings.TrimSpace(noteName)

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
