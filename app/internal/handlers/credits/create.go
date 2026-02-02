package credits

import (
    "context"

	"api/internal/handlers"
	"api/internal/contracts/credits"
	"api/internal/services"
)

func init() {
    handlers.Register(credits.Create, create)
}

func create(ctx context.Context, data map[string]any) (interface{}, error) {
    return services.CreditService.Create(ctx,
        data["client_id"].(int),
        data["bank_id"].(int),
        data["min_payment"].(float64),
        data["max_payment"].(float64),
        data["term_months"].(int),
        data["credit_type"].(string),
    )
}