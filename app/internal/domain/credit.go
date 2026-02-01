package domain

import "time"

type Credit struct {
	ID         string    `json:"id"`
	ClientID   string    `json:"client_id"`
	BankID     string    `json:"bank_id"`
	MinPayment float64   `json:"min_payment"`
	MaxPayment float64   `json:"max_payment"`
	TermMonths int       `json:"term_months"`
	CreditType string    `json:"credit_type"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}