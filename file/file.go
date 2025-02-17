package file

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/rhysmah/note-app/internal/logger"
)

const dateTimeRegexPattern = `(\d{4})_(\d{2})_(\d{2})_(\d{2})_(\d{2})`

var dateTimeRegex = regexp.MustCompile(dateTimeRegexPattern)

type File struct {
	Name         string
	FilePath     string
	DateCreated  time.Time
	DateModified time.Time
}

func NewFile(fileName, notesDir string, logger *logger.Logger) (*File, error) {

	// Validation
	if logger == nil {
		return nil, fmt.Errorf("logger cannot be nil")
	}
	if fileName == "" || notesDir == "" {
		return nil, fmt.Errorf("fileName and notesDir cannot be empty")
	}

	// File Creation
	newFile := &File{
		Name:     fileName,
		FilePath: filepath.Join(notesDir, fileName),
	}

	dateCreated, err := getDateCreated(newFile.FilePath, logger)
	if err != nil {
		return nil, fmt.Errorf("error accessing file's Date Created: %w", err)
	}
	newFile.DateCreated = dateCreated

	dateModified, err := getDateModified(newFile.FilePath, logger)
	if err != nil {
		return nil, fmt.Errorf("error accessing file's Date Modified: %w", err)
	}
	newFile.DateModified = dateModified

	return newFile, nil
}

func getDateModified(filePath string, logger *logger.Logger) (time.Time, error) {
	logger.Start("Getting date modified from file...")

	fileInfo, err := os.Stat(filePath)

	if err != nil {
		errMsg := fmt.Sprintf("error accessing file info: %v", err)
		logger.Fail(errMsg)
		return time.Time{}, fmt.Errorf(errMsg)
	}

	return fileInfo.ModTime(), nil
}

func getDateCreated(filePath string, logger *logger.Logger) (time.Time, error) {

	logger.Start(fmt.Sprintf("Extracting creation date from file %q", filePath))

	matches := dateTimeRegex.FindStringSubmatch(filePath)
	if len(matches) != 6 { // Original + 5 capture groups
		return time.Time{}, fmt.Errorf("invalid filename format: %q", filePath)
	}

	return validateNoteCreatedTime(matches, logger)
}

// HELPERS for GetDateCreated
func validateNoteCreatedTime(dateTimeCmp []string, logger *logger.Logger) (time.Time, error) {
	year, err := strconv.Atoi(dateTimeCmp[1])
	if err != nil {
		logger.Fail(fmt.Sprintf("Failed to convert year %q to integer: %v", dateTimeCmp[1], err))
		return time.Time{}, fmt.Errorf("failed to convert year: %w", err)
	}

	monthAsInt, err := strconv.Atoi(dateTimeCmp[2])
	if err != nil {
		logger.Fail(fmt.Sprintf("Failed to convert month %q to integer: %v", dateTimeCmp[2], err))
		return time.Time{}, fmt.Errorf("failed to convert month: %w", err)
	}
	month := time.Month(monthAsInt)

	day, err := strconv.Atoi(dateTimeCmp[3])
	if err != nil {
		logger.Fail(fmt.Sprintf("Failed to convert day %q to integer: %v", dateTimeCmp[3], err))
		return time.Time{}, fmt.Errorf("failed to convert day: %w", err)
	}

	hour, err := strconv.Atoi(dateTimeCmp[4])
	if err != nil {
		logger.Fail(fmt.Sprintf("Failed to convert hour %q to integer: %v", dateTimeCmp[4], err))
		return time.Time{}, fmt.Errorf("failed to convert hour: %w", err)
	}

	minute, err := strconv.Atoi(dateTimeCmp[5])
	if err != nil {
		logger.Fail(fmt.Sprintf("Failed to convert minute %q to integer: %v", dateTimeCmp[5], err))
		return time.Time{}, fmt.Errorf("failed to convert minute: %w", err)
	}

	creationDate := time.Date(year, month, day, hour, minute, 0, 0, time.Local)

	return creationDate, nil
}

// File Operation Helpers
// ----------------------
// These functions handle the reading and processing of file for the list command.
// They are specific to this command's implementation and shouldn't be used elsewhere.
// ----------------------

// prepareNoteFiles reads and processes notes from the specified directory.
// It returns a slice of File objects or an error if the operation fails.
func PrepareNoteFiles(logger *logger.Logger, notesDir string) ([]File, error) {
	logger.Start(fmt.Sprintf("Preparing notes in directory %q...", notesDir))

	notes, err := ReadNotesDirectory(logger, notesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read notes directory %q: %w", notesDir, err)
	}

	files, err := buildFileObjects(logger, notesDir, notes)
	if err != nil {
		return nil, fmt.Errorf("failed to build File objects for notes in directory %q: %w", notesDir, err)
	}

	return files, nil
}

// buildFileObjects creates File objects from directory entries.
// It processes each note file and returns a slice of File objects.
func buildFileObjects(logger *logger.Logger, notesDir string, notes []os.DirEntry) ([]File, error) {
	files := make([]File, 0, len(notes))

	for _, note := range notes {
		logger.Info(fmt.Sprintf("Processing note: %s", note.Name()))

		newFile, err := NewFile(note.Name(), notesDir, logger)
		if err != nil {
			logger.Fail(fmt.Sprintf("Failed to create File object for %q: %v", note.Name(), err))
			return nil, fmt.Errorf("failed to create File object for %q: %w", note.Name(), err)
		}

		files = append(files, *newFile)
	}

	logger.Success(fmt.Sprintf("Successfully processed %d notes", len(files)))
	return files, nil
}

// readNotesDirectory reads the contents of the notes directory.
// It returns an error if the directory is empty or cannot be read.
func ReadNotesDirectory(logger *logger.Logger, notesDir string) ([]os.DirEntry, error) {
	notes, err := os.ReadDir(notesDir)
	if err != nil {
		logger.Fail(fmt.Sprintf("Failed to read notes directory %q: %v", notesDir, err))
		return nil, fmt.Errorf("failed to read notes directory %q: %w", notesDir, err)
	}

	if len(notes) == 0 {
		logger.Info(fmt.Sprintf("No notes found in directory %q", notesDir))
		return nil, fmt.Errorf("no notes found in directory %q", notesDir)
	}

	logger.Info(fmt.Sprintf("Found %d notes in %q directory", len(notes), notesDir))

	return notes, nil
}
