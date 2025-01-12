package entities

import "time"

type Event struct {
	ID          int           `json:"id" form:"id"`
	UserID      int           `json:"user_id" form:"user_id"`
	Title       string        `json:"title" form:"title"`
	Description string        `json:"description" form:"description"`
	Date        time.Time     `json:"date" form:"date"`
	Duration    time.Duration `json:"duration" form:"duration"`
	CreatedAt   time.Time     `json:"created_at" form:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at" form:"updated_at"`
}
