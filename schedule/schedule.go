package schedule

type EntryID int

type Scheduler interface {
	AddEntry(id EntryID, spec string, cmd func(id EntryID)) error
	RemoveEntry(id EntryID)
	ValidateSpec(spec string) error
}
