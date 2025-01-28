TODOs 

newNote.go
- [ ] Allow user to select location of notes directory. Default: saved as /User/[Username]/notes
- [x] Check that file name is appropriate / correct / won't cause problems
- [x] Add datetime stamp automatically to file name (for searchability, etc.)

log.go
- [x] Create a Logger struct (path to our logger file, pointer to log file)
- [x] Create a New function to create a new Logger
- [x] Create a rotateLogFile function to ensure a previous Log file is closed
- [x] Create a Log function that logs data
	> ensure a log file is open
	> create a timestamp and log entry
	> write data to log file
	> return an error or nil

COMMANDS
- [ ] List - shows the notes a user already has
	> include flags to sort by date
- [ ] View - allows the user to view the contents of the note
	> allow users to view contents in CL or in thei default text editor
- [ ] Edit - allows users to open and edit their notes in default text editor
- [ ] Delete - allows users to delete their notes
	> User selects which note to delete via flag
	> Confirmation occurs so no accidental deletions

list.go
- The os.ReadDir() function automatically sorts files by filename, so
  no specific function needs to be written for that.
- However, it gets complicated if we want to sort for date. There are
  two date sorting functions we can implement: one will sort based on
  Last Modified dates, the other by Date Created dates. However, OS
  files systems only track Last Modified date, so I'll have to add
  the creation date to the filenames themselves, extract that string,
  convert that string into a datetime object, and write a custom
  sorting function using the Go's sort package interface.
- I'll also have to create a File type that stores information, including
  the filename, the creation date, and the last modified date; this is
  required to allow the CLI to work.
- I'll create a method on the File type to extract this data and populate 
  this different file types.

FUTURE FEATURES
- [ ] Add more data for notes -- an object with a name and date field, possibly tags
- [ ] Allow bulk deletion

