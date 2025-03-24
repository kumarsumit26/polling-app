package models

import (
	"time"
)

type Poll struct {
	ID        int       `json:"id"`
	Question  string    `json:"question"`
	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	YesCount  int       `json:"yes_count"`
	NoCount   int       `json:"no_count"`
}

type Vote struct {
	ID        int       `json:"id"`
	PollID    int       `json:"poll_id"`
	Username  string    `json:"username"`
	Vote      bool      `json:"vote"`
	CreatedAt time.Time `json:"created_at"`
}
