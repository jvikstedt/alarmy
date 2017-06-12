package alarmy

import (
	"net/http"
	"time"
)

type Store interface {
	Projects() ([]Project, error)
}

type Router interface {
	Setup() (http.Handler, error)
}

type Project struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
