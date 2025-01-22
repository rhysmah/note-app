package logger

import (
	"fmt"
	"os"
	"path/filepath"
)

// Octal: read/write permission for owner
const ownerReadWritePerms = 0755

type Logger struct {
	currentLogFile *os.File
	logDir string
}

func New() (*Logger, error) {

	// Create app, log dirs if don't exist
	logDir, err := createLogDirectory()
	if err != nil {
		return nil, fmt.Errorf("couldn't create app directory: %w", err)
	}

	// Create intitial log file


	// Set up logger struct (with above data)

	return nil, nil
}

// HELPERS

func createLogDirectory() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("couldn't find user's home directory: %w", err)
	}

	appFilePath := filepath.Join(homeDir, ".note-app")
	logsFilePath := filepath.Join(appFilePath, "logs")

	if os.MkdirAll(logsFilePath, ownerReadWritePerms); err != nil {
		return "", fmt.Errorf("couldn't make app directory: %w", err)
	}
	
	// logsFilePath will be the logDir in Logger struct
	return logsFilePath, nil
}
 