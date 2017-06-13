package alarmy

import (
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

func (s *BoltStore) Close() error {
	return s.db.Close()
}
