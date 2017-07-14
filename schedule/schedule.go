package schedule

type EntryID int

type Scheduler interface {
	AddFunc(id EntryID, spec string, cmd func()) error
	RemoveEntry(id EntryID)
}
