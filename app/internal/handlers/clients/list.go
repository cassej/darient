package clients

import (
    "context"

	"api/internal/handlers"
	"api/internal/contracts/clients"
	"api/internal/services"
)

func init() {
    handlers.Register(clients.List, list)
}

func list(ctx context.Context, data map[string]any) (interface{}, error) {
    page, pageSize := 1, 20

    if v, ok := data["page"].(int); ok && v > 0 {
        page = v
    }
    if v, ok := data["page_size"].(int); ok && v > 0 {
        pageSize = v
    }
    return services.ClientService.List(ctx, page, pageSize)
}