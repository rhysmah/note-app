package cmd

import (
	"fmt"
	"os"

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

		err = listNotes(notesDir)
		if err != nil {
			appLogger.Fail(fmt.Sprintf("Failed to list notes: %v", err))
			fmt.Println(err)
			return
		}
		appLogger.End("Note listing completed successfully")
	},
}

// HELPERS
func listNotes(notesDir string) error {
	appLogger.Start(fmt.Sprintf("Reading notes from directory: %s", notesDir))

	notes, err := os.ReadDir(notesDir)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to read directory: %v", err)
		appLogger.Fail(errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	if len(notes) == 0 {
		appLogger.Info("No notes found in directory")
		fmt.Println("No notes found")
		return nil
	}

	appLogger.Info(fmt.Sprintf("Found %d notes", len(notes)))
	for _, note := range notes {
		fmt.Println(note.Name())
	}

	appLogger.Success("Successfully listed all notes")
	return nil
}
