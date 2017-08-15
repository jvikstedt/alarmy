package model

import "time"

type Result struct {
	ID          int       `json:"id"`
	JobID       int       `json:"job_id"`
	IsTriggered bool      `json:"is_triggered"`
	CreatedAt   time.Time `json:"created_at" db:"createdAt"`
	UpdatedAt   time.Time `json:"updated_at" db:"updatedAt"`
}
