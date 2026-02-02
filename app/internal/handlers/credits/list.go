package credits

import (
    "context"

	"api/internal/handlers"
	"api/internal/contracts/credits"
	"api/internal/services"
)

func init() {
    handlers.Register(credits.List, list)
}

func list(ctx context.Context, data map[string]any) (interface{}, error) {
	page, _ := data["page"].(int)
	pageSize, _ := data["page_size"].(int)

	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 20
	}

	return services.CreditService.List(ctx, page, pageSize)
}