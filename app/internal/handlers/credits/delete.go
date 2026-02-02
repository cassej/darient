package credits

import (
	"api/internal/handlers"
	"api/internal/contracts/credits"
	"api/internal/services"
)

func init() {
    handlers.Register(credits.Delete, delete)
}

func delete(ctx context.Context, data map[string]any) (interface{}, error) {
    return services.CreditService.Delete(ctx, data["id"].(int))
}