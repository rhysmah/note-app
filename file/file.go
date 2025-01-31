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

	// dateCreated, err := extractDateCreated(fileName, logger)
	// if err != nil {
	// 	errMsg := fmt.Sprintf("Error extracting file's Date Created: %v", err)
	// 	logger.Fail(errMsg)
	// 	return nil, fmt.Errorf("%s", errMsg)
	// }

	// newFile.DateCreated = dateCreated

	return newFile, nil
}

func getDateModified(filePath string, logger *logger.Logger) (time.Time, error) {
	logger.Start("Getting date modified from file...")

	fileInfo, err := os.Stat(filePath)

	if err != nil {
		errMsg := fmt.Sprintf("Error accessing file info: %v", err)
		logger.Fail(errMsg)
		return time.Time{}, fmt.Errorf(errMsg)
	}

	return fileInfo.ModTime(), nil
}

// func extractDateCreated(fileName string, logger *logger.Logger) (time.Time, error) {
// 	logger.Start("Getting date created from filename...")

// 	// Find the datetime
// 	regex := regexp.MustCompile(dateTimeRegexPattern)
// 	match := regex.FindStringSubmatch(fileName)

// 	if len(match) != 6 { // Should have full match plus 5 capture groups
// 		errMsg := "Invalid datetime pattern in filename"
// 		logger.Fail(errMsg)
// 		return time.Time{}, fmt.Errorf("%s", errMsg)
// 	}

// 	// Use the captured groups directly instead of splitting again
// 	year, err := strconv.Atoi(match[1])
// 	if err != nil {
// 		errMsg := fmt.Sprintf("Invalid year format: %v", err)
// 		logger.Fail(errMsg)
// 		return time.Time{}, fmt.Errorf("%s", errMsg)
// 	}

// 	month, err := strconv.Atoi(match[2])
// 	if err != nil {
// 		errMsg := fmt.Sprintf("Invalid month format: %v", err)
// 		logger.Fail(errMsg)
// 		return time.Time{}, fmt.Errorf("%s", errMsg)
// 	}

// 	if month < 1 || month > 12 {
// 		errMsg := fmt.Sprintf("Invalid month value: %d", month)
// 		logger.Fail(errMsg)
// 		return time.Time{}, fmt.Errorf("%s", errMsg)
// 	}

// 	day, err := strconv.Atoi(match[3])
// 	if err != nil {
// 		errMsg := fmt.Sprintf("Invalid day format: %v", err)
// 		logger.Fail(errMsg)
// 		return time.Time{}, fmt.Errorf("%s", errMsg)
// 	}
// 	if day < 1 || day > 31 {
// 		errMsg := fmt.Sprintf("Invalid day value: %d", day)
// 		logger.Fail(errMsg)
// 		return time.Time{}, fmt.Errorf("%s", errMsg)
// 	}

// 	hour, err := strconv.Atoi(match[4])
// 	if err != nil {
// 		errMsg := fmt.Sprintf("Invalid hour format: %v", err)
// 		logger.Fail(errMsg)
// 		return time.Time{}, fmt.Errorf("%s", errMsg)
// 	}
// 	if hour < 0 || hour > 23 {
// 		errMsg := fmt.Sprintf("Invalid hour value: %d", hour)
// 		logger.Fail(errMsg)
// 		return time.Time{}, fmt.Errorf("%s", errMsg)
// 	}

// 	minute, err := strconv.Atoi(match[5])
// 	if err != nil {
// 		errMsg := fmt.Sprintf("Invalid minute format: %v", err)
// 		logger.Fail(errMsg)
// 		return time.Time{}, fmt.Errorf("%s", errMsg)
// 	}
// 	if minute < 0 || minute > 59 {
// 		errMsg := fmt.Sprintf("Invalid minute value: %d", minute)
// 		logger.Fail(errMsg)
// 		return time.Time{}, fmt.Errorf("%s", errMsg)
// 	}

// 	creationDateTime := time.Date(
// 		year,
// 		time.Month(month),
// 		day,
// 		hour,
// 		minute,
// 		0, // seconds
// 		0, // nanoseconds
// 		time.Local,
// 	)

// 	logger.Success("Successfully extracted date from filename")
// 	return creationDateTime, nil
// }
