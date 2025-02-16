package list

import (
	"fmt"
	"os"
	"strings"

	"github.com/rhysmah/note-app/cmd/root"
	"github.com/rhysmah/note-app/file"
	"github.com/rhysmah/note-app/internal/filesystem"
	"github.com/rhysmah/note-app/internal/logger"
	"github.com/spf13/cobra"
)

const (
	sortByCmd      = "sort-by"
	sortByCmdShort = "s"

	orderCmd      = "order"
	orderCmdShort = "o"

	listDesc = `List all notes in your notes directory. 
You can sort notes by creation date, modification date, or name.
Example: notes list --sort-by modified --order newest`
)

// init registers the list command and its flags with the root command.
func init() {
	newListCommand := NewListCommand()
	root.RootCmd.AddCommand(newListCommand)

	flags := newListCommand.Flags()

	flags.StringP(sortByCmd, sortByCmdShort, "",
		fmt.Sprintf("Sort by: %s", availableSortFields()))

	flags.StringP(orderCmd, orderCmdShort, string(SortOrderNewest),
		fmt.Sprintf("Order by: %s", availableSortOrders()))
}

// NewListCommand creates and returns a new cobra.Command for the list functionality.
// It handles listing and sorting notes based on user-specified criteria.
func NewListCommand() *cobra.Command {

	listCmd := &ListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List and sort notes",
		Args:  cobra.NoArgs,
		Long:  listDesc,
		RunE: func(cmd *cobra.Command, args []string) error {

			sortBy, err := cmd.Flags().GetString("sort-by")
			if err != nil {
				return fmt.Errorf("failed to get sort-by flag: %w", err)
			}

			order, err := cmd.Flags().GetString("order")
			if err != nil {
				return fmt.Errorf("failed to get order flag: %w", err)
			}

			listCmd.SortField = SortField(sortBy)
			listCmd.SortOrder = SortOrder(order)

			return listCmd.Run(root.AppLogger, root.DirManager)
		},
	}
	return cmd
}

// Run executes the list command with the specified options.
// It completes default values, validates inputs, and processes the notes.
func (opts *ListOptions) Run(logger *logger.Logger, dm *filesystem.DirectoryManager) error {
	if err := opts.Complete(); err != nil {
		return fmt.Errorf("invalid options: %w", err)
	}

	if err := opts.Validate(); err != nil {
		return fmt.Errorf("invalid options: %w", err)
	}

	notesDir := dm.NotesDir()

	logger.Info("Reading notes from directory")
	files, err := prepareNoteFiles(logger, notesDir)
	if err != nil {
		return fmt.Errorf("failed to get files: %w", err)
	}
	opts.files = files

	if err := opts.Execute(); err != nil {
		return fmt.Errorf("could not execute command: %w", err)
	}

	return nil
}

// Complete sets default values for sorting options.
// If no sort field is specified, defaults to sorting by name in alphabetical order.
func (opts *ListOptions) Complete() error {

	if opts.SortField == "" {
		opts.SortField = SortFieldName
		opts.SortOrder = SortOrderAlph
	}

	return nil
}

// Validate checks if the provided options meet all validation rules.
func (opts *ListOptions) Validate() error {
	v := NewValidator()
	return v.Run(opts)
}

// Execute performs the note sorting and displays the results to stdout.
func (opts *ListOptions) Execute() error {

	SortFiles(opts.files, opts.SortField, opts.SortOrder)

	fmt.Println(getHeader(opts.SortField, opts.SortOrder))
	fmt.Println()

	for _, file := range opts.files {
		fmt.Println(file.Name)
	}

	return nil
}

// availableSortFields returns a comma-separated string of valid sort field options.
func availableSortFields() string {
	fields := []string{
		string(SortFieldCreated),
		string(SortFieldModified),
		string(SortFieldName),
	}
	return strings.Join(fields, ", ")
}

// availableSortOrders returns a comma-separated string of valid sort order options.
func availableSortOrders() string {
	orders := []string{
		string(SortOrderNewest),
		string(SortOrderOldest),
		string(SortOrderAlph),
		string(SortOrderRAlph),
	}
	return strings.Join(orders, ", ")
}

// getHeader returns a formatted string describing the current sort configuration.
func getHeader(field SortField, order SortOrder) string {

	fieldDescription := sortFieldDescriptions[field]

	if field == SortFieldName {
		return fmt.Sprintf("Sorting by %s", fieldDescription)
	}
	orderDescription := sortOrderDescriptions[order]

	return fmt.Sprintf("Sorting by %s, %s", fieldDescription, orderDescription)
}


// File Operation Helpers
// ---------------------
// These functions handle the reading and processing of file for the list command. 
// They are specific to this command's implementation and shouldn't be used elsewhere.

// prepareNoteFiles reads and processes notes from the specified directory.
// It returns a slice of File objects or an error if the operation fails.
func prepareNoteFiles(logger *logger.Logger, notesDir string) ([]file.File, error) {
	logger.Start(fmt.Sprintf("Preparing notes in directory %q...", notesDir))

	notes, err := readNotesDirectory(logger, notesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read notes directory %q: %w", notesDir, err)
	}

	files, err := buildFileObjects(logger, notesDir, notes)
	if err != nil {
		return nil, fmt.Errorf("failed to build File objects for notes in directory %q: %w", notesDir, err)
	}

	return files, nil
}

// buildFileObjects creates File objects from directory entries.
// It processes each note file and returns a slice of File objects.
func buildFileObjects(logger *logger.Logger, notesDir string, notes []os.DirEntry) ([]file.File, error) {
	files := make([]file.File, 0, len(notes))

	for _, note := range notes {
		logger.Info(fmt.Sprintf("Processing note: %s", note.Name()))

		newFile, err := file.NewFile(note.Name(), notesDir, logger)
		if err != nil {
			logger.Fail(fmt.Sprintf("Failed to create File object for %q: %v", note.Name(), err))
			return nil, fmt.Errorf("failed to create File object for %q: %w", note.Name(), err)
		}

		files = append(files, *newFile)
	}

	logger.Success(fmt.Sprintf("Successfully processed %d notes", len(files)))
	return files, nil
}

// readNotesDirectory reads the contents of the notes directory.
// It returns an error if the directory is empty or cannot be read.
func readNotesDirectory(logger *logger.Logger, notesDir string) ([]os.DirEntry, error) {
	notes, err := os.ReadDir(notesDir)
	if err != nil {
		logger.Fail(fmt.Sprintf("Failed to read notes directory %q: %v", notesDir, err))
		return nil, fmt.Errorf("failed to read notes directory %q: %w", notesDir, err)
	}

	if len(notes) == 0 {
		logger.Info(fmt.Sprintf("No notes found in directory %q", notesDir))
		return nil, fmt.Errorf("no notes found in directory %q", notesDir)
	}

	logger.Info(fmt.Sprintf("Found %d notes in %q directory", len(notes), notesDir))

	return notes, nil
}
