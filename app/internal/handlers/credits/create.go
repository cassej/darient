package credits

import (
	"api/internal/handlers"
	"api/internal/contracts/credits"
	"api/internal/services"
)

func init() {
    handlers.Register(credits.Create, create)
}

func create(ctx context.Context, data map[string]any) (interface{}, error) {
    return service.CreditService.Create(ctx,
        data["client_id"].(string),
        data["bank_id"].(string),
        data["min_payment"].(float64),
        data["max_payment"].(float64),
        data["term_months"].(int),
        data["credit_type"].(string),
    )
}