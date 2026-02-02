package banks

import (
	"api/internal/handlers"
	"api/internal/contracts/banks"
	"api/internal/services"
)

func init() {
    handlers.Register(banks.Update, update)
}

func update(ctx context.Context, data map[string]any) (interface{}, error) {
    id := data["id"].(string)

    var fullName, email, birthDate, country *string

    if v, ok := data["name"].(string); ok {
        name = &v
    }
    if v, ok := data["type"].(string); ok {
        bank_type = &v
    }

    return service.BankService.Update(ctx, id, name, bank_type)
}