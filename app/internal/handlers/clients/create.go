package clients

import (
	"api/internal/handlers"
	"api/internal/contracts/clients"
	"api/internal/services"
)

func init() {
    handlers.Register(clients.Create, create)
}

func create(ctx context.Context, data map[string]any) (interface{}, error) {
    return service.ClientService.Create(ctx,
        data["full_name"].(string),
        data["email"].(string),
        data["birth_date"].(string),
        data["country"].(string),
    )
}