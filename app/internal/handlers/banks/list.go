package banks

import (
	"net/http"

	"api/internal/handlers"
	"api/internal/middleware"
	"api/internal/repository"
)

func init() {
	handlers.Register("GET", "/banks", ListBanks)
}

func ListBanks(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	pool := middleware.GetDB(r.Context())
	if pool == nil {
		return nil, handlers.NewHTTPError(http.StatusInternalServerError, "database unavailable")
	}

	repo := repository.NewBankRepository(pool)

	banks, err := repo.List(r.Context())
	if err != nil {
		return nil, handlers.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return banks, nil
}