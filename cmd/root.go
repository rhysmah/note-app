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

		dirManager, err = filesystem.NewDirectoryManager(appLogger)
		if err != nil {
			fmt.Printf("Home directory operation failed: %v", err)
			os.Exit(1)
		}
	},

	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if appLogger != nil {
			appLogger.CloseCurrentLogFile()
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
