package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Octal: 4 = read, 2 = write, 1 = execute
const (
	ownerReadWritePerms = 0755
	logReadWritePerms   = 0644
)

const (
	logFilePrefix = "log_"
	logFileSuffix = ".txt"
)

type Logger struct {
	logDirectory   string
	currentLogFile *os.File
	logger         *log.Logger
}

type LogType int

const (
	SuccessLog LogType = iota
	FailLog
	StartLog
	EndLog
	InfoLog
)

func (l LogType) String() string {
	values := [...]string{
		"[SUCCESS]",
		"[FAIL]",
		"[START]",
		"[END]",
		"[INFO]",
	}

	if l < SuccessLog || l > InfoLog {
		return "[UNKNOWN]"
	}

	return values[l]
}

// NewLogger initializes a new Logger instance.
// Creates a new logging directory and the initial log file.
// Returns a pointer to the Logger instance and errors encounted during initialization.
// If any errors occur, it returns nil and the error.
func NewLogger() (*Logger, error) {

	newLogger := &Logger{
		logDirectory:   "",
	}

	// Directory created and assigned to 'logDirectory' field
	logDir, err := newLogger.createLogDirectory()
	if err != nil {
		return nil, fmt.Errorf("couldn't create app directory: %w", err)
	}
	newLogger.logDirectory = logDir

	// Initial log file created and assigned to 'currentLogFile' field
	logFile, err := newLogger.setLoggerFile()
	if err != nil {
		return nil, fmt.Errorf("couldn't create log file: %w", err)
	}
	newLogger.currentLogFile = logFile

	// Leverage built-in 'log' library to display filename and line of error
	newLogger.logger = log.New(newLogger.currentLogFile, "", log.Lshortfile|log.LstdFlags)
	if err := newLogger.Info("Log file initialized"); err != nil {
		newLogger.CloseCurrentLogFile()
		return nil, fmt.Errorf("failed to write initial log entry: %w", err)
	}

	return newLogger, nil
}

func (l *Logger) log(logType LogType, message string) error {

	if l.logDirectory == "" || l.currentLogFile == nil {
		return fmt.Errorf("logger not properly initialized")
	}

	formattedLogMsg := fmt.Sprintf("%s %s\n", logType, message)

	return l.logger.Output(2, formattedLogMsg)
}

// Helpers
func (l *Logger) createLogDirectory() (string, error) {

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("couldn't find user's home directory: %w", err)
	}

	appFilePath := filepath.Join(homeDir, ".note-app")
	logsFilePath := filepath.Join(appFilePath, "logs")

	if err := os.MkdirAll(logsFilePath, ownerReadWritePerms); err != nil {
		return "", fmt.Errorf("couldn't make app directory: %w", err)
	}
	return logsFilePath, nil
}

func (l *Logger) setLoggerFile() (*os.File, error) {

	if l.currentLogFile != nil {
		l.currentLogFile.Close()
	}

	logTimeStamp := time.Now().Format("2006_01_02_15_04")
	logFileName := logFilePrefix + logTimeStamp + logFileSuffix
	logFilePath := filepath.Join(l.logDirectory, logFileName)

	logFile, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, logReadWritePerms)
	if err != nil {
		return nil, fmt.Errorf("error opening log file \"%s\": %w", logFilePath, err)
	}

	return logFile, nil
}

// Writes a log with an "[INFO]" prefix
func (l *Logger) Info(message string) error {
	return l.log(InfoLog, message)
}

// Writes a log with an "[START]" prefix
func (l *Logger) Start(message string) error {
	return l.log(StartLog, message)
}

// Writes a log with an "[END]" prefix
func (l *Logger) End(message string) error {
	return l.log(EndLog, message)
}

// Writes a log with an "[SUCCESS]" prefix
func (l *Logger) Success(message string) error {
	return l.log(SuccessLog, message)
}

// Writes a log with an "[FAIL]" prefix
func (l *Logger) Fail(message string) error {
	return l.log(FailLog, message)
}

// CloseCurrentLogFile checks if the Logger is already using a *os.File and closes it
// If there's a problem closing a file, it returns an error, else returns nil
func (l *Logger) CloseCurrentLogFile() error {
	if l.currentLogFile != nil {
		err := l.currentLogFile.Close()
		l.currentLogFile = nil
		return err
	}

	return nil
}
