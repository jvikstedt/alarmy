package store

import (
	"github.com/jinzhu/gorm"
	"github.com/jvikstedt/alarmy/internal/model"
)

type gormProject struct {
	gormStore *GormStore
}

func (s *gormProject) db() *gorm.DB {
	return s.gormStore.db
}

func (s *gormProject) GetAll(projects interface{}) error {
	return s.db().Find(projects).Error
}

func (s *gormProject) Create(project interface{}) error {
	return s.db().Create(project).Error
}

func (s *gormProject) FindByID(id int, project interface{}) error {
	return s.db().First(project, id).Error
}

func (s *gormProject) CountAll() (int, error) {
	var count int
	err := s.db().Find(&[]*model.Project{}).Count(&count).Error
	return count, err
}

func (s *gormProject) Clear() {
	s.db().Model(&model.Project{}).Delete(&model.Project{})
}
