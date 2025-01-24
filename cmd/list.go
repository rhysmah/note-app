package cmd

import (
	"fmt"
	"os"

	"github.com/rhysmah/note-app/internal/logger"
	"github.com/spf13/cobra"
)

var listLogger *logger.Logger

func init() {
	rootCmd.AddCommand()
}

var list = &cobra.Command{
	Use:   "list",
	Short: "Show all notes",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

// confirmUserHomeDirectory retrieves the user's home directory path from the operating system.
// It returns the absolute path to the user's home directory as a string and any error encountered.
// If the home directory cannot be determined, it returns an empty string and a descriptive error.
func confirmUserHomeDirectory() (string, error) {
	appLogger.Log("[START] Looking up user home directory")

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		errMsg := fmt.Sprintf("[ERROR] Home directory lookup failed: %v", err)
		appLogger.Log(errMsg)
		return "", fmt.Errorf("%s", errMsg)
	}

	appLogger.Log(fmt.Sprintf("[SUCCESS] Found home directory at: %s", userHomeDir))
	return userHomeDir, nil
}
