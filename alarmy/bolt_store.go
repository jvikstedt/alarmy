package alarmy

import (
	"github.com/boltdb/bolt"
)

var BucketKeyProjects = []byte("projects")
var Buckets = [][]byte{
	BucketKeyProjects,
}

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

func (s *BoltStore) EnsureBucketsExist() error {
	return s.CreateBuckets(Buckets...)
}

func (s *BoltStore) RecreateAllBuckets() error {
	return s.RecreateBuckets(Buckets...)
}

func (s *BoltStore) CreateBuckets(names ...[]byte) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		for _, n := range names {
			if _, err := tx.CreateBucketIfNotExists(n); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *BoltStore) RecreateBuckets(names ...[]byte) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		for _, n := range names {
			tx.DeleteBucket(n)
			if _, err := tx.CreateBucket(n); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *BoltStore) Close() error {
	return s.db.Close()
}
