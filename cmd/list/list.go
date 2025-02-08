package list

import (
	"fmt"
	"os"

	"github.com/rhysmah/note-app/cmd/root"
	"github.com/rhysmah/note-app/file"
	"github.com/spf13/cobra"
)

func init() {
	root.RootCmd.AddCommand(list)
	list.PersistentFlags().Bool("c", false, "List all files by their creation date, newest to oldest")
	list.PersistentFlags().Bool("m", false, "List all files by their last-modified date, newest to oldest")
}

var list = &cobra.Command{
	Use:   "list",
	Short: "Show all notes",
	Run: func(cmd *cobra.Command, args []string) {
		root.AppLogger.Start("Listing all notes...")

		byMod, _ := cmd.Flags().GetBool("m")
		byCtd, _ := cmd.Flags().GetBool("c")

		if byMod && byCtd {
			root.AppLogger.Fail("Cannot use both --m and --c flags")
			return
		}

		notesDir := root.DirManager.NotesDir()

		root.AppLogger.Info("Reading notes from directory")
		files, err := getFiles(notesDir)
		if err != nil {
			root.AppLogger.Fail(fmt.Sprintf("Failed to get files: %v", err))
			fmt.Println("Unable to read notes")
			return
		}

		for _, file := range files {
			fmt.Println(file.Name)
		}

		root.AppLogger.End("Note listing completed successfully")
	},
}

func getFiles(notesDir string) ([]file.File, error) {
	notes, err := os.ReadDir(notesDir)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to read directory: %v", err)
		root.AppLogger.Fail(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	if len(notes) == 0 {
		errMsg := "No notes found in directory"
		root.AppLogger.Info(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	root.AppLogger.Info(fmt.Sprintf("Found %d notes", len(notes)))

	files := make([]file.File, 0, len(notes))
	for _, note := range notes {
		root.AppLogger.Info(fmt.Sprintf("Processing note: %s", note.Name()))

		newFile, err := file.NewFile(note.Name(), notesDir, root.AppLogger)

		if err != nil {
			errMsg := fmt.Sprintf("Trouble accessing file: %v", err)
			root.AppLogger.Fail(errMsg)
			return nil, fmt.Errorf(errMsg)
		}

		files = append(files, *newFile)
	}

	root.AppLogger.Success(fmt.Sprintf("Successfully processed %d notes", len(files)))
	return files, nil
}

// TODO:
// 1) Create a sortOptions struct with the booleans of the list options we want.
// 2) Create function to print sortType header
// 3) Create closure function to get sort function typew
// 4) Creat function to validate sort -- only one sort type possible
// 5) Create sortFiles function that sorts files using closure function.
// 6) Create a helper function to ensure the flags we get are bools

// sort.Slice(files, func(i, j int) bool {
// 			switch {
// 			case byMod:
// 				return files[i].DateModified.After(files[j].DateModified)
// 			case byCtd:
// 				return files[i].DateCreated.Before(files[j].DateCreated)
// 			default:
// 				return files[i].Name < files[j].Name
// 			}
// 		})

// 		switch {
// 		case byMod:
// 			fmt.Println("Listing files by Modified Date, newest to oldest")
// 		case byCtd:
// 			fmt.Println("Listing files by Creation Date, newest to oldest")
// 		default:
// 			fmt.Println("Listing files by name")
// 		}
