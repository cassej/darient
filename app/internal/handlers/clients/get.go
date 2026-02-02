package clients

import (
    "context"

	"api/internal/handlers"
	"api/internal/contracts/clients"
	"api/internal/services"
)

func init() {
    handlers.Register(clients.Get, get)
}

func get(ctx context.Context, data map[string]any) (interface{}, error) {
    return services.ClientService.Get(ctx, data["id"].(int))
}