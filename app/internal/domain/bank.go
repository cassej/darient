package domain

import "time"

type Bank struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"` // PRIVATE |   GOVERNMENT
	CreatedAt time.Time `json:"created_at"`
}