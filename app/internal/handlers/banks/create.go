package banks

import (
    "context"

	"api/internal/handlers"
	"api/internal/contracts/banks"
	"api/internal/services"
)

func init() {
    handlers.Register(banks.Create, create)
}

func create(ctx context.Context, data map[string]any) (interface{}, error) {
    return services.BankService.Create(ctx,
        data["name"].(string),
        data["type"].(string),
    )
}