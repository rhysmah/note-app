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

FUTURE FEATURES
- [ ] Add more data for notes -- an object with a name and date field, possibly tags
- [ ] Allow bulk deletion
