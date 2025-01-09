package entities

import "time"

type Event struct {
	ID        int           `json:"id"`
	UserID    int           `json:"user_id"`
	Title     string        `json:"title"`
	Date      time.Time     `json:"date"`
	Duration  time.Duration `json:"duration"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}
