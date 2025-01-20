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

	// `*cobra.Command` is a pointer to the command being executed.
	// `args` are the additional arguments being passed with the command
	// cobra.ExactArgs(1) is a function that returns a function
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(1)(cmd, args); err != nil {
			return fmt.Errorf("you must specify a name when creating a note")
		} 
		return nil // One argument, the name of the note, was provided.
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("A new note was created!")
	},
}