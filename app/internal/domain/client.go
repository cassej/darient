package domain

import "time"

type Client struct {
	ID        int    `json:"id"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	BirthDate string    `json:"birth_date"`
	Country   string    `json:"country"`
	CreatedAt time.Time `json:"created_at"`
}