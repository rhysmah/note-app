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
		appLogger.Log("[START] Listing all notes")

		err := listNotes(dirManager.NotesDir())
		if err != nil {
			appLogger.Log(fmt.Sprintf("[ERROR] Failed to list notes: %v", err))
			fmt.Println(err)
			return
		}
		appLogger.Log("[END] Note listing completed successfully")
	},
}

// HELPERS
func listNotes(notesDir string) error {
	appLogger.Log(fmt.Sprintf("[START] Reading notes from directory: %s", notesDir))

	notes, err := os.ReadDir(notesDir)
	if err != nil {
		errMsg := fmt.Sprintf("[ERROR] Failed to read directory: %v", err)
		appLogger.Log(errMsg)
		return fmt.Errorf(errMsg)
	}

	if len(notes) == 0 {
		appLogger.Log("[INFO] No notes found in directory")
		fmt.Println("No notes found")
		return nil
	}

	appLogger.Log(fmt.Sprintf("[INFO] Found %d notes", len(notes)))
	for _, note := range notes {
		fmt.Println(note.Name())
	}

	appLogger.Log("[SUCCESS] Successfully listed all notes")
	return nil
}
