package model

import "time"

type Project struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	Name      string     `json:"name"`
}

func (p Project) Errors() []Error {
	errors := []Error{}

	if len(p.Name) <= 0 {
		errors = append(errors, Error{Type: "validation", Name: "name", Reason: "missing"})
	}

	return errors
}

type ProjectResponse struct {
	Response
	Project Project `json:"data"`
}

type ProjectListResponse struct {
	Response
	Projects []Project `json:"data"`
}
