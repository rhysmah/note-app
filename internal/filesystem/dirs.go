// Package filesystem provides functionality for managing directory operations in the note-app.
// It handles the creation, validation, and management of essential directories required by the
// application, specifically focusing on the user's home directory and a dedicated notes directory.
//
// The package uses a DirectoryManager type to encapsulate directory-related operations and
// maintain state information about important directory paths. It implements proper error
// handling and logging mechanisms to track operations and handle failure cases appropriately.
//
// Directory permissions are managed using predefined constants, ensuring consistent access
// control across the application. The package relies on the OS's file system operations
// and provides a clean interface for directory management tasks.
package filesystem

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rhysmah/note-app/internal/logger"
)

// TODO: consider enums
const (
	// 7 means read, write, execute for owner
	// 5 means read and execute for group and others, respectively
	dirPermissions  int    = 0755
	defaultNotesDir string = "/notes"
)

// DirectoryManager handles filesystem directory operations for the note-app.
// It manages the application's directory structure and provides access to
// important directory paths like home and notes directories.
// DirectoryManager requires a logger for operation tracking and maintains
// paths for both the home directory and the designated notes directory.
type DirectoryManager struct {
	logger   *logger.Logger
	homeDir  string
	notesDir string
}

// NewDirectoryManager creates and returns a new instance of DirectoryManager with the provided logger.
// The returned DirectoryManager has empty home and notes directory paths, which should be initialized
// using the appropriate setup methods before use.
func NewDirectoryManager(logger *logger.Logger) *DirectoryManager {
	return &DirectoryManager{
		logger:   logger,
		homeDir:  "",
		notesDir: "",
	}
}

// ConfirmUserHomeDirectory retrieves and validates the user's home directory path.
// It uses os.UserHomeDir() to find the home directory and stores it in the DirectoryManager instance.
// The function logs the start of the operation, any errors that occur, and successful completion.
//
// Returns:
//   - string: The path to the user's home directory if successful
//   - error: An error if the home directory lookup fails, nil otherwise
func (dm *DirectoryManager) ConfirmUserHomeDirectory() (string, error) {
	dm.logger.Start("Looking up user home directory")

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		errMsg := fmt.Sprintf("Home directory lookup failed: %v", err)
		dm.logger.Fail(errMsg)
		return "", fmt.Errorf("%s", errMsg)
	}

	dm.logger.Success(fmt.Sprintf("Found home directory at: %s", userHomeDir))
	dm.homeDir = userHomeDir
	return userHomeDir, nil
}

// ConfirmNotesDirectory ensures the existence of the notes directory and returns its path.
// It first validates that the home directory has been set. If not, it returns an error.
// Then it creates the notes directory (including any necessary parent directories) if it doesn't exist.
// The directory is created with the permissions specified by dirPermissions.
//
// Returns:
//   - string: The absolute path to the notes directory
//   - error: An error if the home directory is not set or if directory creation fails
func (dm *DirectoryManager) ConfirmNotesDirectory() (string, error) {
	dm.logger.Start("Setting up notes directory")

	// Check that userHomeDir has been confirmed
	if dm.homeDir == "" {
		errMsg := fmt.Sprintln("User home directory not found")
		dm.logger.Fail(errMsg)
		return "", fmt.Errorf("%s", errMsg)
	}

	notesDirPath := filepath.Join(dm.homeDir, defaultNotesDir)
	dm.logger.Info(fmt.Sprintf("Target notes directory path: %s", notesDirPath))

	err := os.MkdirAll(notesDirPath, os.FileMode(dirPermissions))
	if err != nil {
		errMsg := fmt.Sprintf("Directory creation failed: %v", err)
		dm.logger.Fail(errMsg)
		return "", fmt.Errorf("%s", errMsg)
	}

	dm.logger.Success(fmt.Sprintf("Notes directory ready at: %s", notesDirPath))
	dm.notesDir = notesDirPath
	return notesDirPath, nil
}

func (dm *DirectoryManager) NotesDir() (string, error) {
	if dm.notesDir == "" {
		return "", fmt.Errorf("Notes directory has not been initialized")
	}

	return dm.notesDir, nil
}
