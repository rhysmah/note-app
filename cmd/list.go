package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/rhysmah/note-app/file"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(list)
	list.PersistentFlags().Bool("byCtd", false, "List all files by their creation date, newest to oldest")
	list.PersistentFlags().Bool("byMod", false, "List all files by their last-modified date, newest to oldest")
}

var list = &cobra.Command{
	Use:   "list",
	Short: "Show all notes",
	Run: func(cmd *cobra.Command, args []string) {
		appLogger.Start("Listing all notes...")

		byMod, _ := cmd.Flags().GetBool("byMod")
		// byCtd, _ := cmd.Flags().GetString("byCtd")

		notesDir := dirManager.NotesDir()

		appLogger.Info("Reading notes from directory")
		files, err := getFiles(notesDir)
		if err != nil {
			appLogger.Fail(fmt.Sprintf("Failed to get files: %v", err))
			fmt.Println("Unable to read notes")
			return
		}

		if byMod {
			byModDate := file.ByModifiedDate(files)
			sort.Sort(byModDate)
			for _, file := range byModDate {
				fmt.Println(file.Name)
			}
		} else {
			for _, file := range files {
				fmt.Println(file.Name)
			}
		}
		appLogger.End("Note listing completed successfully")
	},
}

func getFiles(notesDir string) ([]file.File, error) {
	notes, err := os.ReadDir(notesDir)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to read directory: %v", err)
		appLogger.Fail(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	if len(notes) == 0 {
		errMsg := "No notes found in directory"
		appLogger.Info(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	appLogger.Info(fmt.Sprintf("Found %d notes", len(notes)))

	files := make([]file.File, 0, len(notes))
	for _, note := range notes {
		appLogger.Info(fmt.Sprintf("Processing note: %s", note.Name()))

		newFile, err := file.NewFile(note.Name(), notesDir, appLogger)

		if err != nil {
			errMsg := fmt.Sprintf("Trouble accessing file: %v", err)
			appLogger.Fail(errMsg)
			return nil, fmt.Errorf(errMsg)
		}

		files = append(files, *newFile)
	}

	appLogger.Success(fmt.Sprintf("Successfully processed %d notes", len(files)))
	return files, nil
}
