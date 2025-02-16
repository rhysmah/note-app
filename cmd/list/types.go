package list

import "github.com/rhysmah/note-app/file"

type SortField string
type SortOrder string

const (
	SortFieldModified SortField = "modified"
	SortFieldCreated  SortField = "created"
	SortFieldName     SortField = "name"
)

const (
	SortOrderNewest SortOrder = "newest"
	SortOrderOldest SortOrder = "oldest"
	SortOrderAlph   SortOrder = "alphabetical"
	SortOrderRAlph  SortOrder = "reverse alphabetical"
)

// Map SortField types to descriptions
var sortFieldDescriptions = map[SortField]string{
	SortFieldModified: "modification date",
	SortFieldCreated:  "creation date",
	SortFieldName:     "file name",
}

var sortOrderDescriptions = map[SortOrder]string{
	SortOrderNewest: "newest to oldest",
	SortOrderOldest: "oldest to newest",
	SortOrderAlph:   "alphabetical order",
	SortOrderRAlph:  "reverse alphabetical order",
}

type ListOptions struct {
	SortField SortField
	SortOrder SortOrder
	files     []file.File
}
