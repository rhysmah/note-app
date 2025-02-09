package main

import (
	_ "github.com/rhysmah/note-app/cmd/list"
	_ "github.com/rhysmah/note-app/cmd/new"
	"github.com/rhysmah/note-app/cmd/root"
)

func main() {
	root.Execute()
}
