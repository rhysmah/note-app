TODOs 

newNote.go
- [ ] Allow user to select location of notes directory. Default: saved as /User/[Username]/notes
- [x] Check that file name is appropriate / correct / won't cause problems
- [x] Add datetime stamp automatically to file name (for searchability, etc.)

log.go
- [ ] Create a Logger struct (path to our logger file, pointer to log file)
- [ ] Create a New function to create a new Logger
- [ ] Create a rotateLogFile function to ensure a previous Log file is closed
- [ ] Create a Log functiont that logs data
	> ensure a log file is open
	> create a timestamp and log entry
	> write data to log file
	> return an error or nil
