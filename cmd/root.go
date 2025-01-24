/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
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

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "note-app",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

		var err error
		appLogger, err = logger.New()
		if err != nil {
			fmt.Printf("Failed to initialize logger: %v\n", err)
			os.Exit(1)
		}

		dirManager = filesystem.NewDirectoryManager(appLogger)
		_, err = dirManager.ConfirmUserHomeDirectory()
		if err != nil {
			appLogger.Log(fmt.Sprintf("[ERROR] Home directory operation failed: %v", err))
			fmt.Println(err)
			os.Exit(1)
		}

		_, err = dirManager.ConfirmNotesDirectory()
		if err != nil {
			appLogger.Log(fmt.Sprintf("[ERROR] Notes directory creation failed: %v", err))
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

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.note-app.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
