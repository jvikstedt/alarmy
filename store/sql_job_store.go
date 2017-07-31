package store

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/jvikstedt/alarmy/model"
)

type SqlJobStore struct {
	str *SqlStore
}

func NewSqlJobStore(sqlStore *SqlStore) *SqlJobStore {
	return &SqlJobStore{
		str: sqlStore,
	}
}

func (s *SqlJobStore) db() *sqlx.DB {
	return s.str.db
}

func (s *SqlJobStore) All() ([]model.Job, error) {
	jobs := []model.Job{}
	return jobs, s.db().Select(&jobs, "SELECT * FROM jobs")
}

func (s *SqlJobStore) Create(job *model.Job) error {
	result, err := s.db().NamedExec(`
		INSERT INTO jobs (name, projectID, spec, cmd, active, createdAt, updatedAt)
		VALUES (:name, :projectID, :spec, :cmd, :active, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`, job)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	return s.db().Get(job, "SELECT * FROM jobs WHERE id=$1", id)
}

func (s *SqlJobStore) Update(job *model.Job) error {
	result, err := s.db().NamedExec(`
		UPDATE jobs SET (name, projectID, spec, cmd, active, updatedAt)
		= (:name, :projectID, :spec, :cmd, :active, CURRENT_TIMESTAMP) WHERE id=:id`, job)

	if err != nil {
		return err
	}
	id, err := result.RowsAffected()
	if err != nil {
		return err
	}

	return s.db().Get(job, "SELECT * FROM jobs WHERE id=$1", id)
}

func (s *SqlJobStore) Destroy(id int) error {
	result, err := s.db().Exec("DELETE FROM jobs where id=$1", id)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("Could not find record with id of %d", id)
	}

	return nil
}

func (s *SqlJobStore) GetOne(id int) (model.Job, error) {
	job := model.Job{}
	err := s.db().Get(&job, "SELECT * FROM jobs WHERE id=$1", id)
	return job, err
}

func (s *SqlJobStore) RemoveAll() error {
	_, err := s.db().Exec("DELETE FROM jobs")
	return err
}
