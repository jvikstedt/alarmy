package model

import (
	"time"
)

type Project struct {
	ID        int       `json:"id" transform:"lock"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at" transform:"lock"`
	UpdatedAt time.Time `json:"updated_at" transform:"lock"`
}

func (p Project) Errors() map[string][]string {
	errors := make(map[string][]string)

	if len(p.Name) <= 0 {
		errors["name"] = append(errors["name"], "is required")
	}

	return errors
}
