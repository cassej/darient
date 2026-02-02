package events

import(
    "context"
    "time"
)

type EventPublisher interface {
    Publish(ctx context.Context, event Event) error
}

type Event struct {
    Type      string
    Timestamp time.Time
    Payload   any
}
