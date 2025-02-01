package filesystem

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rhysmah/note-app/internal/logger"
)

const (
	// Octal: 4 = read, 2 = write, 1 = execute
	dirPermissions  int    = 0755
	defaultNotesDir string = "/notes"
)

type DirectoryManager struct {
	logger   *logger.Logger
	homeDir  string
	notesDir string
}

func NewDirectoryManager(logger *logger.Logger) *DirectoryManager {
	return &DirectoryManager{
		logger:   logger,
		homeDir:  "",
		notesDir: "",
	}
}

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
