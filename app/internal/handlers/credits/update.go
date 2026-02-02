package credits

import (
	"api/internal/handlers"
	"api/internal/contracts/credits"
	"api/internal/services"
)

func init() {
    handlers.Register(credits.Update, update)
}

func update(ctx context.Context, data map[string]any) (interface{}, error) {
	id := data["id"].(string)

	var minPayment, maxPayment *float64
	var termMonths *int
	var creditType, status *string

	if v, ok := data["min_payment"].(float64); ok {
		minPayment = &v
	}
	if v, ok := data["max_payment"].(float64); ok {
		maxPayment = &v
	}
	if v, ok := data["term_months"].(int); ok {
		termMonths = &v
	}
	if v, ok := data["credit_type"].(string); ok {
		creditType = &v
	}
	if v, ok := data["status"].(string); ok {
		status = &v
	}

	return service.CreditService.Update(ctx, id, minPayment, maxPayment, termMonths, creditType, status)
}