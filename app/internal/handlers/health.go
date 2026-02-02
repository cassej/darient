package handlers

import (
    "context"

	"api/internal/contracts"
)

func init() {
    Register(contracts.Health, health)
}

func health(ctx context.Context, data map[string]any) (interface{}, error) {
    return map[string]string{"status": "ok"}, nil
}