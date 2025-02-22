package delete

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rhysmah/note-app/cmd/root"
	"github.com/spf13/cobra"
)

const (
	delCmd      = "del"
	delCmdShort = "Delete a note"
	delCmdDesc  = `Delete a note with a specified name.
	Usage: note-app del [note-id] or note-app d [note-id]`
)

func init() {
	newDeleteCommand := NewDeleteCommand()
	root.RootCmd.AddCommand(newDeleteCommand)
}

func NewDeleteCommand() *cobra.Command {
	deleteCmd := &DeleteOptions{}

	cmd := &cobra.Command{
		Use:   delCmd,
		Short: delCmdShort,
		Long:  delCmdDesc,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			root.AppLogger.Start(fmt.Sprintf("Deleting note %q", args[0]))

			// NotesDir identified in PersistentPreRun check in root.go
			deleteCmd.notesDir = root.DirManager.NotesDir()
			deleteCmd.noteName = args[0]

			if err := deleteNote(deleteCmd); err != nil {
				errMsg := fmt.Sprintf("Failed to delete note %q", deleteCmd.noteName)
				root.AppLogger.Fail(errMsg)
				return errors.New(errMsg)
			}
			return nil
		},
	}
	return cmd
}

func deleteNote(opts *DeleteOptions) error {
	notePath := filepath.Join(opts.notesDir, opts.noteName)

	if _, err := os.Stat(notePath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("note does not exist %v", err)
		}
		return fmt.Errorf("failed to access note: %w", err)
	}

	if !confirmDeletion(opts.noteName) {
		fmt.Println("User cancelled delete operation")
		return nil
	}

	root.AppLogger.Info(fmt.Sprintf("Deleting note %q...", opts.noteName))

	if err := os.Remove(notePath); err != nil {
		return fmt.Errorf("failed to delete note file: %w", err)
	}

	fmt.Printf("Successfully deleted %q", opts.noteName)
	root.AppLogger.Success(fmt.Sprintf("Note %q successfully deleted", opts.noteName))
	return nil
}

func confirmDeletion(noteName string) bool {
	for {
		var response string

		fmt.Printf("Are you sure you want to delete %q? (y/n): ", noteName)
		fmt.Scanf("%s", &response)

		switch strings.ToLower(response) {
		case "y":
			return true
		case "n":
			return false
		default:
			fmt.Println("Invalid response. Enter (y)es or (n)o.")
		}
	}
}
