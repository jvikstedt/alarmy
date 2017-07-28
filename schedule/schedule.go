package schedule

type EntryID int

type Executor func(EntryID)

type Scheduler interface {
	AddEntry(id EntryID, spec string, executor Executor) error
	RemoveEntry(id EntryID)
	ValidateSpec(spec string) error
}
