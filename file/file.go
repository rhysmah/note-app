package file

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/rhysmah/note-app/internal/logger"
)

const dateTimeRegexPattern = `(\d{4})_(\d{2})_(\d{2})_(\d{2})_(\d{2})$`

type File struct {
	name string
	dateCreated time.Time
	dateModified time.Time
}

func NewFile(fileName string, logger *logger.Logger) (*File, error) {
	newFile := &File{
		name: fileName,
	}

	dateModified, err := getDateModified(fileName, logger)
	if err != nil {
		errMsg := fmt.Sprintf("Error accessing file's Date Modified: %v", err)
		logger.Fail(errMsg)
		return nil, fmt.Errorf("%s", errMsg)
	}
	newFile.dateModified = dateModified

	dateCreated, err := extractDateCreated(fileName, logger)
	if err != nil {
		errMsg := fmt.Sprintf("Error extracting file's Date Created: %v", err)
		logger.Fail(errMsg)
		return nil, fmt.Errorf("%s", errMsg)
	}
	
	newFile.dateCreated = dateCreated

	return newFile, nil
}

// Helper Functions

func getDateModified(fileName string, logger *logger.Logger) (time.Time, error) {
	logger.Start("Getting date modified from file...")

	fileInfo, err := os.Stat(fileName)

	if err != nil {
		errMsg := fmt.Sprintf("Error accessing file info: %v", err)
		logger.Fail(errMsg)

		return time.Time{}, fmt.Errorf("%s", errMsg)
	}

	return fileInfo.ModTime(), nil
}

func extractDateCreated(fileName string, logger *logger.Logger) (time.Time, error) {
    logger.Start("Getting date created from filename...")

    // Find the datetime
    regex := regexp.MustCompile(dateTimeRegexPattern)
    match := regex.FindStringSubmatch(fileName)

    // Check if we found a match at all
    if len(match) < 2 {  // We want at least 2 because match[0] is the full match
        errMsg := "No datetime pattern found in filename"
        logger.Fail(errMsg)
        return time.Time{}, fmt.Errorf("%s", errMsg)
    }

    // Split datetime parts and validate we have enough pieces
    dt := strings.Split(match[1], "_")
    if len(dt) != 5 {
        errMsg := fmt.Sprintf("Invalid datetime format in filename. Expected 5 parts, got %d", len(dt))
        logger.Fail(errMsg)
        return time.Time{}, fmt.Errorf("%s", errMsg)
    }

	year, err := strconv.Atoi(dt[0])
    if err != nil {
        errMsg := fmt.Sprintf("Invalid year format: %v", err)
        logger.Fail(errMsg)
        return time.Time{}, fmt.Errorf("%s", errMsg)
    }

    month, err := strconv.Atoi(dt[1])
    if err != nil {
        errMsg := fmt.Sprintf("Invalid month format: %v", err)
        logger.Fail(errMsg)
        return time.Time{}, fmt.Errorf("%s", errMsg)
    }

    if month < 1 || month > 12 {
        errMsg := fmt.Sprintf("Invalid month value: %d", month)
        logger.Fail(errMsg)
        return time.Time{}, fmt.Errorf("%s", errMsg)
    }

    day, err := strconv.Atoi(dt[2])
    if err != nil {
        errMsg := fmt.Sprintf("Invalid day format: %v", err)
        logger.Fail(errMsg)
        return time.Time{}, fmt.Errorf("%s", errMsg)
    }
    if day < 1 || day > 31 {
        errMsg := fmt.Sprintf("Invalid day value: %d", day)
        logger.Fail(errMsg)
        return time.Time{}, fmt.Errorf("%s", errMsg)
    }

    hour, err := strconv.Atoi(dt[3])
    if err != nil {
        errMsg := fmt.Sprintf("Invalid hour format: %v", err)
        logger.Fail(errMsg)
        return time.Time{}, fmt.Errorf("%s", errMsg)
    }
    if hour < 0 || hour > 23 {
        errMsg := fmt.Sprintf("Invalid hour value: %d", hour)
        logger.Fail(errMsg)
        return time.Time{}, fmt.Errorf("%s", errMsg)
    }

    minute, err := strconv.Atoi(dt[4])
    if err != nil {
        errMsg := fmt.Sprintf("Invalid minute format: %v", err)
        logger.Fail(errMsg)
        return time.Time{}, fmt.Errorf("%s", errMsg)
    }
    if minute < 0 || minute > 59 {
        errMsg := fmt.Sprintf("Invalid minute value: %d", minute)
        logger.Fail(errMsg)
        return time.Time{}, fmt.Errorf("%s", errMsg)
    }

    creationDateTime := time.Date(
        year,
        time.Month(month),
        day,
        hour,
        minute,
        0, // seconds
        0, // nanoseconds
        time.Local,
    )

    logger.Success("Successfully extracted date from filename")
    return creationDateTime, nil 
}