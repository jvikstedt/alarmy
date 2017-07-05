package store

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/boltdb/bolt"
	"github.com/jvikstedt/alarmii/util"
	"github.com/jvikstedt/alarmy/model"
)

type BoltProjectStore struct {
	*BoltStore
}

func NewBoltProjectStore(bs *BoltStore) *BoltProjectStore {
	return &BoltProjectStore{
		BoltStore: bs,
	}
}

func (s *BoltProjectStore) ProjectAll() ([]model.Project, error) {
	projects := []model.Project{}
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(BucketKeyProjects)

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			project := model.Project{}
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

func (s *BoltProjectStore) ProjectCreate(project model.Project) (model.Project, error) {
	err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(BucketKeyProjects)

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

func (s *BoltProjectStore) ProjectUpdate(project model.Project) (model.Project, error) {
	err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(BucketKeyProjects)

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

func (s *BoltProjectStore) ProjectDestroy(id int) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(BucketKeyProjects)

		return b.Delete(util.Itob(id))
	})
}

func (s *BoltProjectStore) ProjectGetOne(id int) (model.Project, error) {
	project := model.Project{}
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(BucketKeyProjects)

		data := b.Get(util.Itob(id))
		if len(data) <= 0 {
			return fmt.Errorf("No record found with id of %d", id)
		}
		return json.Unmarshal(data, &project)
	})
	return project, err
}

func (s *BoltProjectStore) ProjectRemoveAll() error {
	return s.RecreateBuckets(BucketKeyProjects)
}
