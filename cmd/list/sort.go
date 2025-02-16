package list

import (
	"sort"

	"github.com/rhysmah/note-app/file"
)

// sortFiles sorts a slice of files based on the specified field and order.
// It uses the compareFiles function to determine the ordering between any two files.
func sortFiles(files []file.File, field SortField, order SortOrder) {
	sort.Slice(files, func(a, b int) bool {
		return compareFiles(files[a], files[b], field, order)
	})
}

// compareFiles compares two files based on the specified sort field and order.
// It returns true if file 'a' should come before file 'b' in the sorted result.
func compareFiles(a, b file.File, field SortField, order SortOrder) bool {
	switch field {

	case SortFieldName:
		if order == SortOrderAlph {
			return a.Name < b.Name
		}
		return a.Name > b.Name

	case SortFieldCreated:
		if order == SortOrderNewest {
			return a.DateCreated.After(b.DateCreated)
		}
		return a.DateCreated.Before(b.DateCreated)

	case SortFieldModified:
		if order == SortOrderNewest {
			return a.DateModified.After(b.DateModified)
		}
		return a.DateCreated.Before(b.DateModified)

	default:
		return a.Name < b.Name
	}
}
