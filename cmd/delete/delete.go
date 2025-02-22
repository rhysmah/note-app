package delete

import (
	"github.com/rhysmah/note-app/cmd/root"
	"github.com/spf13/cobra"
)

const (
	delCmd 	   = "del"
	delCmdShort = "d"
	delCmdDesc  = `Delete a note with a specified name.
	Usage: note-app del [note-id] or note-app d [note-id]
`
)

func NewDeleteCommand() *cobra.Command {
	deleteCmd := &DeleteOptions{}

	cmd := &cobra.Command{
		Use: delCmd,
		Short: delCmdShort,
		Long: delCmdDesc,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			// NotesDir identified in PersistentPreRun check in root.go
			deleteCmd.notesDir = root.DirManager.NotesDir()
			deleteCmd.noteName = args[0]

			// Delete command runs here


			return nil
		},
	}
	
	return cmd
}