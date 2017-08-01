package store

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/jvikstedt/alarmy/model"
)

type SqlProjectStore struct {
	str *SqlStore
}

func NewSqlProjectStore(bs *SqlStore) *SqlProjectStore {
	return &SqlProjectStore{
		str: bs,
	}
}

func (s *SqlProjectStore) db() *sqlx.DB {
	return s.str.db
}

func (s *SqlProjectStore) All() ([]model.Project, error) {
	projects := []model.Project{}
	return projects, s.db().Select(&projects, "SELECT * FROM projects")
}

func (s *SqlProjectStore) Create(project *model.Project) error {
	result, err := s.db().NamedExec(`
		INSERT INTO projects (name, createdAt, updatedAt)
		VALUES (:name, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`, project)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	return s.db().Get(project, "SELECT * FROM projects WHERE id=$1", id)
}

func (s *SqlProjectStore) Update(project *model.Project) error {
	result, err := s.db().NamedExec(`
		UPDATE projects SET (name, updatedAt)
		= (:name, CURRENT_TIMESTAMP) WHERE id=:id`, project)

	if err != nil {
		return err
	}
	id, err := result.RowsAffected()
	if err != nil {
		return err
	}

	return s.db().Get(project, "SELECT * FROM projects WHERE id=$1", id)
}

func (s *SqlProjectStore) Destroy(id int) error {
	result, err := s.db().Exec("DELETE FROM projects where id=$1", id)
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

func (s *SqlProjectStore) GetOne(id int) (model.Project, error) {
	project := model.Project{}
	err := s.db().Get(&project, "SELECT * FROM projects WHERE id=$1", id)
	return project, err
}

func (s *SqlProjectStore) RemoveAll() error {
	_, err := s.db().Exec("DELETE FROM projects")
	return err
}
