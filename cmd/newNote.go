package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(newNote)
}

var newNote = &cobra.Command{
	Use: "create",
	Short: "Create a new note",
	Long: "Create a new note inside the [] directory with the name [note-name]_[date].txt",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("A new note was created!")
	},
}