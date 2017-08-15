package model

import "time"

type TriggerType int

const (
	TriggerEqual TriggerType = iota
	TriggerLessThan
	TriggerMoreThan
)

type Trigger struct {
	ID        uint        `gorm:"primary_key" json:"id"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	DeletedAt *time.Time  `json:"deleted_at"`
	Target    string      `json:"target"`
	Val       string      `json:"val"`
	Type      TriggerType `json:"type"`
	Job       Job         `json:"job"`
	JobID     int         `json:"job_id"`
}

func (t Trigger) Errors() map[string][]string {
	errors := make(map[string][]string)

	if t.JobID <= 0 {
		errors["job_id"] = append(errors["job_id"], "is required")
	}

	if len(t.Target) <= 0 {
		errors["target"] = append(errors["target"], "is required")
	}

	return errors
}
