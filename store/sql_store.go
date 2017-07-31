package store

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

var schema = `
CREATE TABLE IF NOT EXISTS projects (
	id        INTEGER    PRIMARY KEY,
	name      CHAR(100)  NOT NULL,
	createdAt DATETIME   DEFAULT  CURRENT_TIMESTAMP,
	updatedAt DATETIME   DEFAULT  CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS jobs (
	id        INTEGER  PRIMARY KEY,
	name      CHAR(100)  NOT NULL,
	projectID INTEGER,
	spec      CHAR(100),
	cmd       TEXT,
	active    BOOLEAN,
	createdAt DATETIME  DEFAULT  CURRENT_TIMESTAMP,
	updatedAt DATETIME  DEFAULT  CURRENT_TIMESTAMP,

	FOREIGN KEY(projectID) REFERENCES projects(id)
);
`

type SqlStore struct {
	db      *sqlx.DB
	project ProjectStore
	job     JobStore
}

func NewSqlStore(db *sql.DB, driverName string) *SqlStore {
	store := &SqlStore{
		db: sqlx.NewDb(db, driverName),
	}

	store.project = NewSqlProjectStore(store)
	store.job = NewSqlJobStore(store)

	return store
}

func (s *SqlStore) SetupTables() error {
	_, err := s.db.Exec(schema)
	return err
}

func (s *SqlStore) Project() ProjectStore {
	return s.project
}

func (s *SqlStore) Job() JobStore {
	return s.job
}
