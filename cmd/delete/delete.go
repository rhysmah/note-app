package delete

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rhysmah/note-app/cmd/root"
	"github.com/spf13/cobra"
)

const (
	delCmd      = "del"
	delCmdShort = "d"
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

			// NotesDir identified in PersistentPreRun check in root.go
			deleteCmd.notesDir = root.DirManager.NotesDir()
			deleteCmd.noteName = args[0]

			if err := deleteNote(deleteCmd); err != nil {
				fmt.Printf("Error deleting note %q", deleteCmd.noteName)
				return err
			}

			fmt.Printf("Note %q successfully deleted", deleteCmd.noteName)
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

	root.AppLogger.Info(fmt.Sprintf("Deleting note %q...", opts.noteName))

	if err := os.Remove(notePath); err != nil {
		return fmt.Errorf("failed to delete note file: %w", err)
	}

	root.AppLogger.Success(fmt.Sprintf("Note %q successfully deleted", opts.noteName))
	return nil
}
