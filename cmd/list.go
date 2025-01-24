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
		appLogger.Log("[START]: Starting to list notes")

		_, err := dirManager.ConfirmUserHomeDirectory()
		if err != nil {
			appLogger.Log(fmt.Sprintf("[ERROR] Home directory operation failed: %v", err))
			fmt.Println(err)
			return
		}

		notesDir, err := dirManager.ConfirmNotesDirectory()
		if err != nil {
			appLogger.Log(fmt.Sprintf("[ERROR] Notes directory creation failed: %v", err))
			fmt.Println(err)
			return
		}

		err = listNotes(notesDir)
		if err != nil {
			fmt.Println(err)
		}
	},
}

// HELPERS
func listNotes(notesDir string) error {
	notes, err := os.ReadDir(notesDir)
	if err != nil {
		return err
	}

	for _, note := range notes {
		fmt.Println(note.Name())
	}
	return nil
}
