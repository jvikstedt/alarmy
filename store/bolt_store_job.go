package store

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/boltdb/bolt"
	"github.com/jvikstedt/alarmii/util"
	"github.com/jvikstedt/alarmy/model"
)

type BoltJobStore struct {
	*BoltStore
}

func NewBoltJobStore(bs *BoltStore) *BoltJobStore {
	return &BoltJobStore{
		BoltStore: bs,
	}
}

func (s *BoltJobStore) JobAll() ([]model.Job, error) {
	projects := []model.Job{}
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(BucketKeyJobs)

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			project := model.Job{}
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

func (s *BoltJobStore) JobCreate(project model.Job) (model.Job, error) {
	err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(BucketKeyJobs)

		id, err := b.NextSequence()
		if err != nil {
			return err
		}
		project.ID = int(id)
		project.CreatedAt = time.Now()
		project.UpdatedAt = time.Now()
		encoded, err := json.Marshal(project)
		if err != nil {
			return err
		}
		return b.Put(util.Itob(project.ID), encoded)
	})
	return project, err
}

func (s *BoltJobStore) JobUpdate(project model.Job) (model.Job, error) {
	err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(BucketKeyJobs)

		if project.ID == 0 {
			return fmt.Errorf("Can't update project with id of 0")
		}

		project.UpdatedAt = time.Now()
		encoded, err := json.Marshal(project)
		if err != nil {
			return err
		}
		return b.Put(util.Itob(project.ID), encoded)
	})
	return project, err
}

func (s *BoltJobStore) JobDestroy(id int) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(BucketKeyJobs)

		return b.Delete(util.Itob(id))
	})
}

func (s *BoltJobStore) JobGetOne(id int) (model.Job, error) {
	project := model.Job{}
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(BucketKeyJobs)

		data := b.Get(util.Itob(id))
		if len(data) <= 0 {
			return fmt.Errorf("No record found with id of %d", id)
		}
		return json.Unmarshal(data, &project)
	})
	return project, err
}

func (s *BoltJobStore) JobRemoveAll() error {
	return s.RecreateBuckets(BucketKeyJobs)
}
