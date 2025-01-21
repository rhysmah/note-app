package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// Octal, as indiciated by 0.
// 7 means read, write, execute for owner
// 5 means read and execture for group and others, respectively
const dirPermissions int = 0755
const defaultNotesDir string = "/notes"
const illegalChars string = "\\/:*?\"<>|:."
const noteNameCharLimit = 50

func init() {
	rootCmd.AddCommand(newNote)
}

var newNote = &cobra.Command{
	Use: "create",
	Short: "Create a new note",
	Long: "Create a new note inside the [] directory with the name [note-name]_[date].txt",

	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(1)(cmd, args); err != nil {
			return fmt.Errorf("you must specify a name when creating a note")
		} 
		return nil
	},

	Run: func(cmd *cobra.Command, args []string) {
		err := validateNoteName(args[0])
		if err != nil {
			fmt.Println(err)
			return
		}
		
		userHomeDir, err := confirmUserHomeDirectory()
		if err != nil {
			fmt.Println(err)
			return
		}
		notesDir, err := createNotesDirectory(userHomeDir)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = createAndSaveNote(notesDir, args[0])
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}

// HELPER FUNCTIONS

// confirmUserHomeDirectory retrieves the user's home directory path from the operating system.
// It returns the absolute path to the user's home directory as a string and any error encountered.
// If the home directory cannot be determined, it returns an empty string and a descriptive error.
func confirmUserHomeDirectory() (string, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("couldn't find user's home directory: %v", err)
	}
	return userHomeDir, nil
}

// createNotesDirectory creates a directory for storing notes in the user's home directory.
// It takes the user's home directory path as input and returns the full path to the created notes directory.
// If creation fails, it returns an empty string and an error describing what went wrong.
// The directory is created with permissions specified by dirPermissions constant.
func createNotesDirectory(userHomeDir string) (string, error) {
	notesDirPath := filepath.Join(userHomeDir, defaultNotesDir)

	err := os.MkdirAll(notesDirPath, os.FileMode(dirPermissions))
	if err != nil {
		return "", fmt.Errorf("couldn't make notes directory: %v", err)
	}
	return notesDirPath, nil
}

// createAndSaveNote creates a new text file with the given name in the specified directory.
// It takes two parameters:
//   - notesDirPath: the directory path where the note should be created
//   - noteName: the name of the note (without .txt extension)
//
// The function will append ".txt" to the note name automatically.
// It returns an error if:
//   - a note with the same name already exists
//   - there are problems creating the file
//
// On success, it prints the path where the note was created and returns nil.
func createAndSaveNote(notesDirPath, noteName string) error {
	
	fullNoteName := noteName + ".txt"
	notePath := filepath.Join(notesDirPath, fullNoteName) 

	// Check if note already exist (do not overwrite)
	if _, err := os.Stat(notePath); err == nil {
		return fmt.Errorf("%s already exists", noteName)
	}

	// Create new note if it does not exist
	file, err := os.Create(notePath)
	if err != nil {
		return fmt.Errorf("problem creating note: %v", err)
	}
	defer file.Close()

	fmt.Printf("Note successfully created at %s\n", notePath)
	return nil
}

// validateNoteName checks if the provided note name meets the required criteria.
// It verifies that:
// - The note name does not exceed the character limit
// - The note name does not begin or end with whitespace
// - The note name does not contain any illegal characters (., \, /, :, *, ?, ", <, >, |)
//
// Parameters:
//   - noteName: string to be validated
//
// Returns:
//   - error: nil if validation passes, error with description if validation fails
func validateNoteName(noteName string) error {
	if len(noteName) > noteNameCharLimit {
		return fmt.Errorf("note name cannot exceed %d characters", noteNameCharLimit)
	}

	// Note cannot start or end with whitespace
	if strings.TrimSpace(noteName) != noteName {
		return fmt.Errorf("note name cannot begin or end with spaces")
	}
	
	// Note cannot contain illegal characters
	var illegalCharsFound []rune
	for _, r := range noteName {
		if strings.ContainsRune(illegalChars, r) {
			illegalCharsFound = append(illegalCharsFound, r)
		}
	}
	
	if len(illegalCharsFound) > 0 {
		return fmt.Errorf("note name contains illegal character(s): \"%s\"", string(illegalCharsFound))
	}

	return nil
}


