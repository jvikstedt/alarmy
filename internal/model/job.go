package model

import "time"

type Job struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	Name      string     `json:"name"`
	Project   Project    `json:"project"`
	ProjectID int        `json:"project_id"`
	Spec      string     `json:"spec"`
	Cmd       string     `json:"cmd"`
	Active    bool       `json:"active"`
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
