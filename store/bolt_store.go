package store

import (
	"github.com/boltdb/bolt"
)

var BucketKeyProjects = []byte("projects")
var BucketKeyJobs = []byte("jobs")
var Buckets = [][]byte{
	BucketKeyProjects,
	BucketKeyJobs,
}

type BoltStore struct {
	db      *bolt.DB
	project ProjectStore
	job     JobStore
}

func NewBoltStore(filename string) (*BoltStore, error) {
	db, err := bolt.Open(filename, 0600, nil)
	if err != nil {
		return nil, err
	}

	store := &BoltStore{
		db: db,
	}
	store.project = NewBoltProjectStore(store)
	store.job = NewBoltJobStore(store)

	err = store.EnsureTablesExist()
	if err != nil {
		return nil, err
	}

	return store, nil
}

func (s *BoltStore) Project() ProjectStore {
	return s.project
}

func (s *BoltStore) Job() JobStore {
	return s.job
}

func (s *BoltStore) EnsureTablesExist() error {
	return s.CreateBuckets(Buckets...)
}

func (s *BoltStore) RecreateAllTables() error {
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
