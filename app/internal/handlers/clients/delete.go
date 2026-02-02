package clients

import (
	"api/internal/handlers"
	"api/internal/contracts/clients"
	"api/internal/services"
)

func init() {
    handlers.Register(clients.Delete, delete)
}

func delete(ctx context.Context, data map[string]any) (interface{}, error) {
    return services.ClientService.Delete(ctx, data["id"].(int))
}