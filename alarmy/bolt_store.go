package alarmy

import "fmt"

type BoltStore struct {
}

func NewBoltStore() *BoltStore {
	return &BoltStore{}
}

func (s *BoltStore) Projects() ([]Project, error) {
	return []Project{}, fmt.Errorf("Not yet implemented")
}
