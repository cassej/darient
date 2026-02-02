package events

import (
    "context"
    "encoding/json"
    "github.com/redis/go-redis/v9"
)

type RedisPublisher struct {
    client *redis.Client
}

func NewRedisPublisher(client *redis.Client) *RedisPublisher {
    return &RedisPublisher{client: client}
}

func (p *RedisPublisher) Publish(ctx context.Context, event Event) error {
    data, _ := json.Marshal(event)
    return p.client.XAdd(ctx, &redis.XAddArgs{
        Stream: "credit_events",
        Values: map[string]any{
            "type":    event.Type,
            "payload": data,
        },
    }).Err()
}