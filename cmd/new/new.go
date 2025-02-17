package new

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/rhysmah/note-app/cmd/root"
	"github.com/rhysmah/note-app/file"
	"github.com/spf13/cobra"
)

const (
	illegalChars      string = "\\/:*?\"<>|: ."
	noteNameCharLimit int    = 50
)

const (
	createCmd 	   = "create"
	createCmdShort = "c"
	createCmdDesc  = `Create a new note with the specified name.
The note will be saved as '[note-name]_[date].txt' in your notes directory.
Note names cannot contain special characters or exceed 50 characters.`
)

func init() {
	newCreateCommand := NewCreateCommand()
	root.RootCmd.AddCommand(newCreateCommand)
}

func NewCreateCommand() *cobra.Command {
	newCmd := &NewOptions{}

	cmd := &cobra.Command{
	Use:   createCmd,
	Short: createCmdShort,
	Long:  createCmdDesc,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		root.AppLogger.Start(fmt.Sprintf("Creating new note with name: '%s'", args[0]))

		dir, err := file.ReadNotesDirectory(root.AppLogger, root.DirManager.HomeDir())
		if err != nil {
			return fmt.Errorf("error: %w", err)
		}
		newCmd.noteDir = dir.name


		root.AppLogger.End("Note creation process completed successfully")
		},
	}
	return cmd
}







func createNote(noteName string) error {

	if err := validateNoteName(noteName); err != nil {
		return fmt.Errorf("invalid note name: %w", err)
	}

	if err := createAndSaveNote(noteName, notesDir); err != nil {
		return fmt.Errorf("failed to create note %s: %w", noteName, err)
	}

	return nil
}

func createAndSaveNote(noteName, notesDirPath string) error {
	root.AppLogger.Start(fmt.Sprintf("Creating note '%s' in directory %s...", noteName, notesDirPath))

	fullNoteName := noteName + "_" + time.Now().Format("2006_01_02_15_04") + ".txt"
	notePath := filepath.Join(notesDirPath, fullNoteName)

	// Check if note already exists
	if _, err := os.Stat(notePath); err == nil {
		errMsg := fmt.Sprintf("note %q already exists", fullNoteName)
		root.AppLogger.Fail(errMsg)
		return fmt.Errorf(errMsg)
	}

	// Create note
	file, err := os.Create(notePath)
	if err != nil {
		errMsg := fmt.Sprintf("failed to create file: %v", err)
		root.AppLogger.Fail(errMsg)
		return fmt.Errorf(errMsg)
	}
	defer file.Close()

	successMsg := fmt.Sprintf("note created at: %s", notePath)
	root.AppLogger.Success(successMsg)
	fmt.Printf("Created note: %s\n", notePath)
	return nil
}
