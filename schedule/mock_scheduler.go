package schedule

type MockScheduler struct {
}

func (s *MockScheduler) AddEntry(id EntryID, spec string, executor Executor) error {
	return nil
}
func (s *MockScheduler) RemoveEntry(id EntryID) {
}
func (s *MockScheduler) ValidateSpec(spec string) error {
	return nil
}
