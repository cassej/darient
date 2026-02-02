package clients

import (
	"api/internal/handlers"
	"api/internal/contracts/clients"
	"api/internal/services"
)

func init() {
    handlers.Register(clients.Update, update)
}

func update(ctx context.Context, data map[string]any) (interface{}, error) {
    id := data["id"].(string)

    var fullName, email, birthDate, country *string

    if v, ok := data["full_name"].(string); ok {
        fullName = &v
    }
    if v, ok := data["email"].(string); ok {
        email = &v
    }
    if v, ok := data["birth_date"].(string); ok {
        birthDate = &v
    }
    if v, ok := data["country"].(string); ok {
        country = &v
    }

    return service.ClientService.Update(ctx, id, fullName, email, birthDate, country)
}