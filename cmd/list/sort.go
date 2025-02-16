package list

import (
	"sort"

	"github.com/rhysmah/note-app/file"
)

func SortFiles(files []file.File, field SortField, order SortOrder) {
	sort.Slice(files, func(a, b int) bool {
		return compareFiles(files[a], files[b], field, order)
	})
}

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
