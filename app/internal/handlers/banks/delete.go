package banks

import (
	"api/internal/handlers"
	"api/internal/contracts/banks"
	"api/internal/services"
)

func init() {
    handlers.Register(banks.Delete, delete)
}

func delete(ctx context.Context, data map[string]any) (interface{}, error) {
    return services.BankService.Delete(ctx, data["id"].(int))
}