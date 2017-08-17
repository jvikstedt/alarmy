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

func (s *gormProject) GetAll(projects *[]model.Project) error {
	return s.db().Find(projects).Error
}

func (s *gormProject) Create(project *model.Project) error {
	return s.db().Create(project).Error
}

func (s *gormProject) Find(project *model.Project) error {
	return s.db().First(project, project.ID).Error
}

func (s *gormProject) Delete(project *model.Project) error {
	err := s.db().Where("id = ?", project.ID).Delete(project).Error
	if err != nil {
		return err
	}
	return s.db().Unscoped().First(project, project.ID).Error
}

func (s *gormProject) Update(project *model.Project) error {
	return s.db().Save(project).Error
}

func (s *gormProject) CountAll() (int, error) {
	var count int
	err := s.db().Find(&[]model.Project{}).Count(&count).Error
	return count, err
}

func (s *gormProject) Clear() {
	s.db().Model(&model.Project{}).Delete(&model.Project{})
}
