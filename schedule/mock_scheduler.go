package schedule

type MockScheduler struct {
}

func (s *MockScheduler) AddFunc(id EntryID, spec string, cmd func(id EntryID)) error {
	return nil
}
func (s *MockScheduler) RemoveEntry(id EntryID) {
}
func (s *MockScheduler) ValidateSpec(spec string) error {
	return nil
}
