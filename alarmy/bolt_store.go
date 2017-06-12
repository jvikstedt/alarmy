package alarmy

import (
	"fmt"

	"github.com/boltdb/bolt"
)

type BoltStore struct {
	db *bolt.DB
}

func NewBoltStore(filename string) (*BoltStore, error) {
	db, err := bolt.Open(filename, 0600, nil)
	if err != nil {
		return nil, err
	}
	return &BoltStore{
		db: db,
	}, nil
}

func (s *BoltStore) Projects() ([]Project, error) {
	return []Project{}, fmt.Errorf("Not yet implemented")
}

func (s *BoltStore) Close() error {
	return s.db.Close()
}
