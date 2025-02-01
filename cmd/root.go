package cmd

import (
	"fmt"
	"os"

	"github.com/rhysmah/note-app/internal/filesystem"
	"github.com/rhysmah/note-app/internal/logger"
	"github.com/spf13/cobra"
)

var (
	appLogger  *logger.Logger
	dirManager *filesystem.DirectoryManager
)

var rootCmd = &cobra.Command{
	Use: "note-app",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

		var err error
		appLogger, err = logger.NewLogger()
		if err != nil {
			fmt.Printf("Failed to initialize logger: %v\n", err)
			os.Exit(1)
		}

		dirManager = filesystem.NewDirectoryManager(appLogger)
		_, err = dirManager.ConfirmUserHomeDirectory()
		if err != nil {
			appLogger.Fail(fmt.Sprintf("Home directory operation failed: %v", err))
			fmt.Println(err)
			os.Exit(1)
		}

		_, err = dirManager.ConfirmNotesDirectory()
		if err != nil {
			appLogger.Fail(fmt.Sprintf("Notes directory creation failed: %v", err))
			fmt.Println(err)
			os.Exit(1)
		}
	},

	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if appLogger != nil {
			appLogger.Close()
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
