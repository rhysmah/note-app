package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Octal: read/write/execute permission for owner
// 4 = read (r), 2 = write (w), 1 = execute (x)
const ownerReadWritePerms = 0755
const logReadWritePerms = 0644
const logFilePrefix = "note_app_log_"

// Logger manages logging operations with file rotation capabilities.
// It maintains a reference to the current log file and the directory
// where log files are stored.
type Logger struct {
	currentLogFile *os.File
	logDir         string
}

// New creates and initializes a new Logger instance.
// It sets up the necessary directory structure for logging and creates an initial log file.
// Returns a pointer to the Logger instance and any error encountered during initialization.
// If an error occurs during directory creation or log file setup, it returns nil and the error.
func New() (*Logger, error) {

	logger := &Logger{}

	err := logger.createLogDirectory()
	if err != nil {
		return nil, fmt.Errorf("couldn't create app directory: %w", err)
	}

	if err := logger.setLoggerFile(); err != nil {
		return nil, fmt.Errorf("couldn't create log file: %w", err)
	}

	return logger, nil
}

func (l *Logger) Log(message string) error {

	if l.currentLogFile == nil {
		return fmt.Errorf("no log file is currently open")
	}

	messageTimeStamp := time.Now().Format("2006-01-02 15:04:05")
	logEntry := fmt.Sprintf("[%s] %s\n", messageTimeStamp, message)

	_, err := l.currentLogFile.WriteString(logEntry)
	if err != nil {
		return fmt.Errorf("error writing to log file %s: %w", l.currentLogFile.Name(), err)
	}

	return nil
}

func (l *Logger) Close() error {
	if l.currentLogFile != nil {
		return l.currentLogFile.Close()
	}
	return nil
}

// HELPERS

// createLogDirectory creates the necessary directory structure for application logs.
// It creates a .note-app directory in the user's home directory, followed by a logs subdirectory.
// The function returns the full path to the logs directory and any error encountered during creation.
// The created directories will have owner read-write permissions set.
// Returns:
//   - string: The absolute path to the created logs directory
//   - error: An error if directory creation fails or if home directory cannot be determined
func (l *Logger) createLogDirectory() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("couldn't find user's home directory: %w", err)
	}

	appFilePath := filepath.Join(homeDir, ".note-app")
	logsFilePath := filepath.Join(appFilePath, "logs")

	if err := os.MkdirAll(logsFilePath, ownerReadWritePerms); err != nil {
		return fmt.Errorf("couldn't make app directory: %w", err)
	}

	l.logDir = logsFilePath
	return nil
}

// setLoggerFile creates and sets up a new log file for the Logger.
// If there's an existing log file, it will be closed before creating a new one.
// The new log file is created with a timestamp in its name following the pattern:
// "[prefix]_YYYY_MM_DD_HH_mm.txt".
// The file is opened in append mode with read-write permissions.
// Returns an error if the file cannot be created or opened.
func (l *Logger) setLoggerFile() error {
	// Ensure previous log file closed
	if l.currentLogFile != nil {
		l.currentLogFile.Close()
	}

	// Create new log file with timestamp
	logTimeStamp := time.Now().Format("2006_01_02_15_04")
	logFileName := filepath.Join(l.logDir, logFilePrefix + logTimeStamp + ".txt")

	// Create log file
	logFile, err := os.OpenFile(logFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, logReadWritePerms)
	if err != nil {
		return fmt.Errorf("error opening log file \"%s\": %w", logFileName, err)
	}

	l.currentLogFile = logFile

	if err := l.Log("Log file initialized"); err != nil {
		return fmt.Errorf("failed to writing initial log entry: %w", err)
	}
	return nil
}
