package file

type ByCreationDate []File

func (bcd ByCreationDate) Len() int { return len(bcd) }
func (bcd ByCreationDate) Swap(i, j int) { bcd[i], bcd[j] = bcd[j], bcd[i] }
func (bcd ByCreationDate) Less(i, j int) bool { return bcd[i].DateCreated.Before(bcd[j].DateCreated) }

type ByModifiedDate []File

func (bmd ByModifiedDate) Len() int { return len(bmd) }
func (bmd ByModifiedDate) Swap(i, j int) { bmd[i], bmd[j] = bmd[j], bmd[i] }
func (bmd ByModifiedDate) Less(i, j int) bool { return bmd[i].DateCreated.Before(bmd[j].DateCreated) }
