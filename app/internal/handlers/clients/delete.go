package clients

import (
    "context"

	"api/internal/handlers"
	"api/internal/contracts/clients"
	"api/internal/services"
)

func init() {
    handlers.Register(clients.Delete, delete)
}

func delete(ctx context.Context, data map[string]any) (interface{}, error) {
    id, _ := data["id"].(int)

    if err := services.ClientService.Delete(ctx, data["id"].(int)); err != nil {
        return nil, err
    }

    return map[string]any{
        "status":  "success",
        "id":      id,
    }, nil
}