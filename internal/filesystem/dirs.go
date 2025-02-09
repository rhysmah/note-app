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

func NewDirectoryManager(logger *logger.Logger) (*DirectoryManager, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger cannot be nil")
	}

	dm := &DirectoryManager{
		logger: logger,
	}

	if err := dm.initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize directory manager: %w", err)
	}

	return dm, nil
}

func (dm *DirectoryManager) initialize() error {
	homeDir, err := dm.confirmUserHomeDirectory()
	if err != nil {
		return fmt.Errorf("failed to confirm home directory: %w", err)
	}
	dm.homeDir = homeDir

	notesDir, err := dm.confirmNotesDirectory()
	if err != nil {
		return fmt.Errorf("failed to confirm notes directory: %w", err)
	}
	dm.notesDir = notesDir

	return nil
}

func (dm *DirectoryManager) HomeDir() string {
	return dm.homeDir
}

func (dm *DirectoryManager) NotesDir() string {
	return dm.notesDir
}

func (dm *DirectoryManager) confirmUserHomeDirectory() (string, error) {
	dm.logger.Start("Looking up user home directory...")

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		errMsg := fmt.Sprintf("Home directory lookup failed: %v", err)
		dm.logger.Fail(errMsg)
		return "", fmt.Errorf("%s", errMsg)
	}

	dm.logger.Success(fmt.Sprintf("Found home directory at: %s", userHomeDir))
	return userHomeDir, nil
}

func (dm *DirectoryManager) confirmNotesDirectory() (string, error) {
	dm.logger.Start("Setting up notes directory...")

	notesDirPath := filepath.Join(dm.homeDir, defaultNotesDir)
	dm.logger.Info(fmt.Sprintf("Target notes directory path: %s", notesDirPath))

	err := os.MkdirAll(notesDirPath, os.FileMode(dirPermissions))
	if err != nil {
		errMsg := fmt.Sprintf("Directory creation failed: %v", err)
		dm.logger.Fail(errMsg)
		return "", fmt.Errorf("%s", errMsg)
	}

	dm.logger.Success(fmt.Sprintf("Notes directory ready at: %s", notesDirPath))
	return notesDirPath, nil
}
