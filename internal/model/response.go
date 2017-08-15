package model

type Error struct {
	Type   string `json:"type"`
	Name   string `json:"name"`
	Reason string `json:"reason"`
}

type Response struct {
	HTTPStatusCode int         `json:"-"`
	Data           interface{} `json:"data"`
	HasErrors      bool        `json:"has_errors"`
	Errors         []Error     `json:"errors"`
}
