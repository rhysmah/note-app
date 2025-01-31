package cmd

import (
	"fmt"
	"os"

	"github.com/rhysmah/note-app/file"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(list)
}

var list = &cobra.Command{
	Use:   "list",
	Short: "Show all notes",
	Run: func(cmd *cobra.Command, args []string) {
		appLogger.Start("Listing all notes")

		notesDir, err := dirManager.NotesDir()
		if err != nil {
			appLogger.Fail(fmt.Sprintf("Cannot access notes directory: %v", err))
			fmt.Println(err)
			return
		}

		appLogger.Info("Reading notes from directory")
		files, err := getFiles(notesDir)
		if err != nil {
			appLogger.Fail(fmt.Sprintf("Failed to get files: %v", err))
			fmt.Println("Unable to read notes")
			return
		}

		for _, file := range files {
			fmt.Println(file.Name)
		}

		appLogger.End("Note listing completed successfully")
	},
}

// HELPERS
func getFiles(notesDir string) ([]file.File, error) {

	notes, err := os.ReadDir(notesDir)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to read directory: %v", err)
		appLogger.Fail(errMsg)
		return nil, fmt.Errorf("%s", errMsg)
	}

	if len(notes) == 0 {
		appLogger.Info("No notes found in directory")
		errMsg := "No notes found in directory"
		return nil, fmt.Errorf("%s", errMsg)
	}

	appLogger.Info(fmt.Sprintf("Found %d notes", len(notes)))
	appLogger.Info(fmt.Sprintf("Starting to process %d notes", len(notes)))

	files := make([]file.File, 0, len(notes))

	for _, note := range notes {
		appLogger.Info(fmt.Sprintf("Processing note: %s", note.Name()))

		newFile, err := file.NewFile(note.Name(), notesDir, appLogger)
		if err != nil {
			errMsg := fmt.Sprintf("Trouble accessing file: %v", err)
			appLogger.Fail(errMsg)
			return nil, fmt.Errorf("%s", errMsg)
		}

		files = append(files, *newFile)
	}

	appLogger.Success(fmt.Sprintf("Successfully processed %d notes", len(files)))
	return files, nil
}
