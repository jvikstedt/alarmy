package model

import "time"

const (
	TriggerEqual = iota
	TriggerLessThan
	TriggerMoreThan
)

type TriggerType uint8

type Trigger struct {
	FieldName string
	Target    string
	TriggerType
}

type Job struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	ProjectID int       `json:"project_id"`
	Spec      string    `json:"spec"`
	Cmd       string    `json:"cmd"`
	Active    bool      `json:"active"`
	Triggers  []Trigger `json:"triggers"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
