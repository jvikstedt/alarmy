package model

import "time"

type Job struct {
	ID        int       `json:"id" transform:"lock"`
	Name      string    `json:"name"`
	ProjectID int       `json:"project_id"`
	Spec      string    `json:"spec"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at" transform:"lock"`
	UpdatedAt time.Time `json:"updated_at" transform:"lock"`
}

func (j Job) Errors() map[string][]string {
	errors := make(map[string][]string)

	if len(j.Name) <= 0 {
		errors["name"] = append(errors["name"], "is required")
	}

	if j.ProjectID == 0 {
		errors["project_id"] = append(errors["project_id"], "is required")
	}

	return errors
}
