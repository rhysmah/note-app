package root

import (
	"fmt"
	"os"

	"github.com/rhysmah/note-app/internal/filesystem"
	"github.com/rhysmah/note-app/internal/logger"
	"github.com/spf13/cobra"
)

var (
	AppLogger      *logger.Logger
	DirManager     *filesystem.DirectoryManager
	UserDirectory  string
	NotesDirectory string
)

var RootCmd = &cobra.Command{
	Use: "note-app",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

		var err error
		AppLogger, err = logger.NewLogger()
		if err != nil {
			fmt.Printf("Failed to initialize logger: %v\n", err)
			os.Exit(1)
		}

		DirManager, err = filesystem.NewDirectoryManager(AppLogger)
		if err != nil {
			fmt.Printf("Failed to initialize logger: %v", err)
			os.Exit(1)
		}
	},

	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if AppLogger != nil {
			AppLogger.CloseCurrentLogFile()
		}
	},
}

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
