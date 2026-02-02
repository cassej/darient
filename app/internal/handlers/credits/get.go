package credits

import (
	"api/internal/handlers"
	"api/internal/contracts/credits"
	"api/internal/services"
)

func init() {
    handlers.Register(credits.Get, get)
}

func get(ctx context.Context, data map[string]any) (interface{}, error) {
    return services.CreditService.Get(ctx, data["id"].(int))
    )
}