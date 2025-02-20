package list

import (
	"fmt"
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

	flags.StringP(orderCmd, orderCmdShort, "",
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
	if err := opts.complete(); err != nil {
		return fmt.Errorf("invalid options: %w", err)
	}

	if err := opts.validate(); err != nil {
		return fmt.Errorf("invalid options: %w", err)
	}

	notesDir := dm.NotesDir()

	logger.Info("Reading notes from directory")
	files, err := file.PrepareNoteFiles(logger, notesDir)
	if err != nil {
		return fmt.Errorf("failed to get files: %w", err)
	}
	opts.files = files

	if err := opts.execute(); err != nil {
		return fmt.Errorf("could not execute command: %w", err)
	}

	return nil
}

// complete sets default values for sorting options.
// If no sort field is specified, defaults to sorting by name in alphabetical order.
func (opts *ListOptions) complete() error {
	if opts.SortField == "" {
		opts.SortField = SortFieldName
		opts.SortOrder = SortOrderAlph
	}

	return nil
}

// validate checks if the provided options meet all validation rules.
func (opts *ListOptions) validate() error {
	v := NewValidator()
	return v.Run(opts)
}

// execute performs the note sorting and displays the results to stdout.
func (opts *ListOptions) execute() error {
	sortFiles(opts.files, opts.SortField, opts.SortOrder)

	//TODO: Improve how files are displayed
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
