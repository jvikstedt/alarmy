package store

import "github.com/jinzhu/gorm"

type GormStore struct {
	db      *gorm.DB
	project ProjectStore
}

func NewGormStore(db *gorm.DB) *GormStore {
	store := &GormStore{
		db: db,
	}

	store.project = &gormProject{gormStore: store}

	return store
}

func (s *GormStore) Project() ProjectStore {
	return s.project
}
