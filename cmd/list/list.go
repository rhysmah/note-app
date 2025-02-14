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

func init() {
	newListCommand := NewListCommand()
	root.RootCmd.AddCommand(newListCommand)

	newListCommand.Flags().String("sort-by", "",
		fmt.Sprintf("Sort by: %s", availableSortFields()))
	newListCommand.Flags().String("order", string(SortOrderNewest),
		fmt.Sprintf("Sort order: %s", availableSortOrders()))
}

// This will return a validated, ready-to-run list command
func NewListCommand() *cobra.Command {

	// This has methods on it that will validate and populate the fields before it's returned.
	listCmd := &ListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List and sort notes",
		Long: `List all notes in your notes directory. 
You can sort notes by creation date, modification date, or name.
Example: notes list --sort-by modified --order newest`,
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

// Set defaults: we want to ensure that if not options are selected that
// defaults ARE selected, so no errors occur.
func (opts *ListOptions) Complete() error {

	if opts.SortField == "" {
		opts.SortField = SortFieldName
	}

	if opts.SortOrder == "" && opts.SortField != SortFieldName {
		opts.SortOrder = SortOrderNewest
	}

	return nil
}

// We want to validate the options selected. For example, there are
// certain flags that only work with certain commands, so we have
// to ensure those are caught and returned.
func (opts *ListOptions) Validate() error {

	// Check Sort Field -- we pass in the command provided by the user
	// to our map; and if the user provides an invalid command, we let them know
	if _, valid := sortFieldDescriptions[opts.SortField]; !valid {
		return fmt.Errorf("invalid sort field selected: %q. Valid sort fields: %q", opts.SortField, availableSortFields())
	}

	// Only validate order if we're not sorting by name
	if opts.SortField != SortFieldName {
		if _, valid := sortOrderDescriptions[opts.SortOrder]; !valid {
			return fmt.Errorf("invalid sort order selected: %q. Valid sort orders: %q",
				opts.SortOrder, availableSortOrders())
		}
	}

	return nil
}

// Does the actual sorting; this assumes that everything is
// working; it has no knowledge of the previous functions,
// of the validation.
func (opts *ListOptions) Execute() error {

	SortFiles(opts.files, opts.SortField, opts.SortOrder)

	fmt.Println(getHeader(opts.SortField, opts.SortOrder))
	fmt.Println()

	for _, file := range opts.files {
		fmt.Println(file.Name)
	}

	return nil
}

func availableSortFields() string {
	fields := []string{
		string(SortFieldCreated),
		string(SortFieldModified),
		string(SortFieldName),
	}
	return strings.Join(fields, ", ")
}

func availableSortOrders() string {
	orders := []string{
		string(SortOrderNewest),
		string(SortOrderOldest),
	}
	return strings.Join(orders, ", ")
}

func getHeader(field SortField, order SortOrder) string {
	fieldDescription := sortFieldDescriptions[field]
	if field == SortFieldName {
		return fmt.Sprintf("Sorting by %s", fieldDescription)
	}
	orderDescription := sortOrderDescriptions[order]
	return fmt.Sprintf("Sorting by %s, %s", fieldDescription, orderDescription)
}

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

// HELPER for getFiles()
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

// HELPER for getFiles()
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
