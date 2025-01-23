// Package cmd provides command-line functionality for managing notes,
// including creation and organization of text files in a specified directory.s
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rhysmah/note-app/internal/logger"
	"github.com/spf13/cobra"
)

const (
	// Octal, as indiciated by 0.
	// 7 means read, write, execute for owner
	// 5 means read and execute for group and others, respectively
	dirPermissions    int    = 0755
	defaultNotesDir   string = "/notes"
	illegalChars      string = "\\/:*?\"<>|:."
	noteNameCharLimit int    = 50
)

var appLogger *logger.Logger

// init initializes the command structure by adding the newNote command
// as a subcommand to the root command. This function is automatically
// called during package initialization.
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
		appLogger.Log(fmt.Sprintf("[START] Creating new note with name: '%s'", args[0]))

		err := validateNoteName(args[0])
		if err != nil {
			appLogger.Log(fmt.Sprintf("[ERROR] Name validation failed for '%s': %v", args[0], err))
			fmt.Println(err)
			return
		}
		appLogger.Log(fmt.Sprintf("[SUCCESS] Name validation passed for '%s'", args[0]))

		userHomeDir, err := confirmUserHomeDirectory()
		if err != nil {
			appLogger.Log(fmt.Sprintf("[ERROR] Home directory operation failed: %v", err))
			fmt.Println(err)
			return
		}

		notesDir, err := createNotesDirectory(userHomeDir)
		if err != nil {
			appLogger.Log(fmt.Sprintf("[ERROR] Notes directory creation failed: %v", err))
			fmt.Println(err)
			return
		}

		err = createAndSaveNote(notesDir, args[0])
		if err != nil {
			appLogger.Log(fmt.Sprintf("[ERROR] Note creation failed: %v", err))
			fmt.Println(err)
			return
		}
		appLogger.Log("[END] Note creation process completed successfully")
	},
}

// HELPER FUNCTIONS

// confirmUserHomeDirectory retrieves the user's home directory path from the operating system.
// It returns the absolute path to the user's home directory as a string and any error encountered.
// If the home directory cannot be determined, it returns an empty string and a descriptive error.
func confirmUserHomeDirectory() (string, error) {
	appLogger.Log("[START] Looking up user home directory")

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		errMsg := fmt.Sprintf("[ERROR] Home directory lookup failed: %v", err)
		appLogger.Log(errMsg)
		return "", fmt.Errorf("%s", errMsg)
	}

	appLogger.Log(fmt.Sprintf("[SUCCESS] Found home directory at: %s", userHomeDir))
	return userHomeDir, nil
}

// createNotesDirectory creates a directory for storing notes in the user's home directory.
// It takes the user's home directory path as input and returns the full path to the created notes directory.
// If creation fails, it returns an empty string and an error describing what went wrong.
// The directory is created with permissions specified by dirPermissions constant.
func createNotesDirectory(userHomeDir string) (string, error) {
	appLogger.Log("[START] Setting up notes directory")

	notesDirPath := filepath.Join(userHomeDir, defaultNotesDir)
	appLogger.Log(fmt.Sprintf("[INFO] Target notes directory path: %s", notesDirPath))

	err := os.MkdirAll(notesDirPath, os.FileMode(dirPermissions))
	if err != nil {
		errMsg := fmt.Sprintf("[ERROR] Directory creation failed: %v", err)
		appLogger.Log(errMsg)
		return "", fmt.Errorf("%s", errMsg)
	}

	appLogger.Log(fmt.Sprintf("[SUCCESS] Notes directory ready at: %s", notesDirPath))
	return notesDirPath, nil
}

// createAndSaveNote creates a new text file with the given name and current timestamp
// in the specified directory. The file name format is: noteName_YYYY_MM_DD_HH_mm.txt
//
// Parameters:
//   - notesDirPath: the directory path where the note will be saved
//   - noteName: the base name for the note (without extension)
//
// Returns:
//   - error: nil if successful, otherwise returns an error if:
//   - the note already exists
//   - there are problems creating the file
func createAndSaveNote(notesDirPath, noteName string) error {
	appLogger.Log(fmt.Sprintf("[START] Creating note file for: '%s'", noteName))

	fullNoteName := noteName + "_" + time.Now().Format("2006_01_02_15_04") + ".txt"
	notePathFinal := filepath.Join(notesDirPath, fullNoteName)

	appLogger.Log(fmt.Sprintf("[INFO] Target note path: %s", notePathFinal))

	if _, err := os.Stat(notePathFinal); err == nil {
		errMsg := fmt.Sprintf("[ERROR] Note already exists: %s", fullNoteName)
		appLogger.Log(errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	file, err := os.Create(notePathFinal)
	if err != nil {
		errMsg := fmt.Sprintf("[ERROR] File creation failed: %v", err)
		appLogger.Log(errMsg)
		return fmt.Errorf("problem creating note: %v", err)
	}
	defer file.Close()

	successMsg := fmt.Sprintf("[SUCCESS] Note created at: %s", notePathFinal)
	appLogger.Log(successMsg)
	fmt.Printf("Note successfully created at %s\n", notePathFinal)
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
	appLogger.Log(fmt.Sprintf("[START] Validating note name: '%s'", noteName))

	if len(noteName) > noteNameCharLimit {
		errMsg := fmt.Sprintf("[ERROR] Name exceeds %d character limit", noteNameCharLimit)
		appLogger.Log(errMsg)
		return fmt.Errorf("note name cannot exceed %d characters", noteNameCharLimit)
	}

	if strings.TrimSpace(noteName) != noteName {
		errMsg := "[ERROR] Name contains leading or trailing spaces"
		appLogger.Log(errMsg)
		return fmt.Errorf("note name cannot begin or end with spaces")
	}

	if strings.Contains(illegalChars, noteName) {
		errMsg := "[ERROR] Name contains illegal characters"
		appLogger.Log(errMsg)
		return fmt.Errorf("note name cannot contain any of the following illegal characters:\\, /, ., *, ?, \", <, >, :, |")
	}

	appLogger.Log("[SUCCESS] Note name validation passed all checks")
	return nil
}
