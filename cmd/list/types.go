package list

import "github.com/rhysmah/note-app/file"

// Example: we may sort the FIELD Date Modified in the ORDER newest to oldest
type SortField string  // The field used to sort files
type SortOrder string  // The order to sort the field

const (
	SortFieldModified SortField = "modified"
	SortFieldCreated SortField = "created"
	SortFieldName SortField = "name"
)

const (
	SortOrderNewest SortOrder = "newest"
	SortOrderOldest SortOrder = "oldest"
)

// Map SortField types to descriptions
var sortFieldDescriptions = map[SortField]string {
	SortFieldModified: "modification date",
	SortFieldCreated: "creation date",
	SortFieldName: "file name",
}

var sortOrderDescriptions = map[SortOrder]string {
	SortOrderNewest: "newest to oldest",
	SortOrderOldest: "oldest to newest",
}

type ListOptions struct {
	SortField SortField
	SortOrder SortOrder
	files []file.File
}