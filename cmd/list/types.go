package list

import "github.com/rhysmah/note-app/file"

type SortField string
type SortOrder string

const (
	SortFieldModified SortField = "mod"
	SortFieldCreated  SortField = "ctd"
	SortFieldName     SortField = "name"
)

const (
	SortOrderNewest SortOrder = "new"
	SortOrderOldest SortOrder = "old"
	SortOrderAlph   SortOrder = "alph"
	SortOrderRAlph  SortOrder = "ralph"
)

const (
	SortFieldModifiedDesc = "modification date"
	SortFieldCreatedDesc  = "creation date"
	SortFieldNameDesc     = "file name"

	SortOrderNewestDesc = "newest to oldest"
	SortOrderOldestDesc = "oldest to newest"
	SortOrderAlphDesc   = "alphabetical"
	SortOrderRAlphDesc  = "reverse alphabetical"
)

var sortFieldDescriptions = map[SortField]string{
	SortFieldModified: SortFieldModifiedDesc,
	SortFieldCreated:  SortFieldCreatedDesc,
	SortFieldName:     SortFieldNameDesc,
}

var sortOrderDescriptions = map[SortOrder]string{
	SortOrderNewest: SortOrderNewestDesc,
	SortOrderOldest: SortOrderOldestDesc,
	SortOrderAlph:   SortOrderAlphDesc,
	SortOrderRAlph:  SortOrderRAlphDesc,
}

type ListOptions struct {
	SortField SortField
	SortOrder SortOrder
	files     []file.File
}
