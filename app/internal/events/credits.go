package events

import "time"

type CreditCreatedEvent struct {
    CreditID   int
    ClientID   int
    BankID     int
    Amount     float64
    CreditType string
}

type CreditApprovedEvent struct {
    CreditID   int
    ClientID   int
    ApprovedAt time.Time
}