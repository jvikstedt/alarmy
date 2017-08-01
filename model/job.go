package model

import (
	"time"
)

type Job struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	ProjectID int       `json:"project_id" db:"projectID"`
	Spec      string    `json:"spec"`
	Cmd       string    `json:"cmd"`
	Active    bool      `json:"active"`
	Triggers  []Trigger `json:"triggers"`
	CreatedAt time.Time `json:"created_at" db:"createdAt"`
	UpdatedAt time.Time `json:"updated_at" db:"updatedAt"`
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

type JobParams struct {
	Job
	OmitID        interface{} `json:"id,omitempty"`
	OmitCreatedAt interface{} `json:"created_at,omitempty"`
	OmitUpdatedAt interface{} `json:"updated_at,omitempty"`
}
