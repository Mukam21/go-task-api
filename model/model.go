package model

import "time"

type Task struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Status    string    `json:"status"` // new, in_progress, done
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
