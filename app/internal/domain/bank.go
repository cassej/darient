package domain

import "time"
import "errors"

var (
	ErrInvalidInput = errors.New("invalid input parameters")
	ErrNotFound     = errors.New("resource not found")
)

type Bank struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"` // PRIVATE |   GOVERNMENT
	CreatedAt time.Time `json:"created_at"`
}