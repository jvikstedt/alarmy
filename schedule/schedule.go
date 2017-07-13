package schedule

type EntryID int

type Scheduler interface {
	AddFunc(spec string, cmd func()) (EntryID, error)
	Remove(id EntryID)
}
