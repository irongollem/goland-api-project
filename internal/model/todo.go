package model

type TodoItem struct {
	ID     int    `json:"id,omitempty"`
	Task   string `json:"task"`
	Status bool   `json:"status"`
}
