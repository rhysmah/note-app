package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// Octal, as indiciated by 0.
// 7 means read, write, execute for owner
// 5 means read and execture for group and others, respectively
const dirPermissions int = 0755

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
		userHomeDir, err := confirmUserHomeDirectory()
		if err != nil {
			fmt.Println(err)
			return
		}
		
		notesDir, err := createNotesDirectory(userHomeDir)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = createAndSaveNote(notesDir, args[0])
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}

// TODO: Allow user to select location of notes directory
// For now, it will be saved as /User/[Username]/notes

// HELPER FUNCTIONS
func confirmUserHomeDirectory() (string, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("couldn't find user's home directory: %v", err)
	}
	return userHomeDir, nil
}

func createNotesDirectory(userHomeDir string) (string, error) {
	var notesDir = "/notes"

	notesDirPath := filepath.Join(userHomeDir, notesDir)

	err := os.MkdirAll(notesDirPath, os.FileMode(dirPermissions))
	if err != nil {
		return "", fmt.Errorf("couldn't make notes directory: %v", err)
	}
	return notesDirPath, nil
}

// createAndSaveNote creates a new text file with the given name in the specified directory.
// It takes two parameters:
//   - notesDirPath: the directory path where the note should be created
//   - noteName: the name of the note (without .txt extension)
//
// The function will append ".txt" to the note name automatically.
// It returns an error if:
//   - a note with the same name already exists
//   - there are problems creating the file
//
// On success, it prints the path where the note was created and returns nil.
func createAndSaveNote(notesDirPath, noteName string) error {
	
	fullNoteName := noteName + ".txt"
	notePath := filepath.Join(notesDirPath, fullNoteName) 

	// Check if note already exist (do not overwrite)
	if _, err := os.Stat(fullNoteName); err == nil {
		return fmt.Errorf("note %s already exists", noteName)
	}

	// Create new note if it does not exist
	file, err := os.Create(notePath)
	if err != nil {
		return fmt.Errorf("problem creating note: %v", err)
	}
	defer file.Close()

	fmt.Printf("Note was created at %s", notePath)
	return nil
}
