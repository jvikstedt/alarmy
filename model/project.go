package model

import (
	"time"
)

type Project struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	ProjectIDs []int     `json:"project_ids"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (p Project) Errors() map[string][]string {
	errors := make(map[string][]string)

	if len(p.Name) <= 0 {
		errors["name"] = append(errors["name"], "is required")
	}

	return errors
}
