package alarmy

import (
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
)

var bucketProjects = []byte("projects")

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

func (s *BoltStore) CreateBucketsIfNotExists() error {
	return s.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucketProjects)
		if err != nil {
			return err
		}
		return nil
	})
}

func (s *BoltStore) Projects() ([]Project, error) {
	projects := []Project{}
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketProjects)

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			project := Project{}
			err := json.Unmarshal(v, &project)
			if err != nil {
				return err
			}
			projects = append(projects, project)
		}
		return nil
	})
	return projects, err
}

func (s *BoltStore) Close() error {
	return s.db.Close()
}
