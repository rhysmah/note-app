package file

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/rhysmah/note-app/internal/logger"
)

// const dateTimeRegexPattern = `(\d{4})_(\d{2})_(\d{2})_(\d{2})_(\d{2})$`

type File struct {
	Name         string
	FilePath     string
	DateCreated  time.Time
	DateModified time.Time
}

func NewFile(fileName, notesDir string, logger *logger.Logger) (*File, error) {

	newFile := &File{
		Name:     fileName,
		FilePath: filepath.Join(notesDir, fileName),
	}

	dateModified, err := getDateModified(newFile.FilePath, logger)
	if err != nil {
		errMsg := fmt.Sprintf("Error accessing file's Date Modified: %v", err)
		logger.Fail(errMsg)
		return nil, fmt.Errorf("%s", errMsg)
	}
	newFile.DateModified = dateModified

	return newFile, nil
}

func getDateModified(filePath string, logger *logger.Logger) (time.Time, error) {
	logger.Start("Getting date modified from file...")

	fileInfo, err := os.Stat(filePath)

	if err != nil {
		errMsg := fmt.Sprintf("Error accessing file info: %v", err)
		logger.Fail(errMsg)
		return time.Time{}, fmt.Errorf("%s", errMsg)
	}

	return fileInfo.ModTime(), nil
}
