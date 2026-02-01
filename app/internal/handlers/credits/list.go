package credits

import (
	"net/http"
	"strconv"

	"api/internal/handlers"
	"api/internal/middleware"
	"api/internal/repository"
	baseRepo "api/pkg/repository"
)

func init() {
	handlers.Register("GET", "/credits", ListCredits)
}

func ListCredits(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	pagination := baseRepo.NewPaginationParams(page, pageSize)

	pool := middleware.GetDB(r.Context())
	if pool == nil {
		return nil, handlers.NewHTTPError(http.StatusInternalServerError, "database unavailable")
	}

	repo := repository.NewCreditRepository(pool)

	result, err := repo.List(r.Context(), pagination)
	if err != nil {
		return nil, handlers.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return result, nil
}