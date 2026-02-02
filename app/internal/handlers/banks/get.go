package banks

import (
	"api/internal/handlers"
	"api/internal/contracts/banks"
	"api/internal/services"
)

func init() {
    handlers.Register(banks.Get, get)
}

func get(ctx context.Context, data map[string]any) (interface{}, error) {
    return services.BankService.Get(ctx, data["id"].(int))
    )
}