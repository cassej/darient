package credits

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"api/internal/contracts"
	"api/internal/domain"
	"api/internal/handlers"
	"api/internal/middleware"
	"api/internal/repository"
)

func init() {
	handlers.Register("GET", "/credits/{id}", GetCredit)
}

func GetCredit(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	id := chi.URLParam(r, "id")
	if id == "" {
		return nil, domain.ErrInvalidInput
	}

    err := contracts.ValidateField("id", id, contracts.FieldSpec{
        Type: "uuid",
    })

    if err != nil {
        return nil, handlers.NewHTTPError(http.StatusBadRequest, err.Error())
    }

	pool := middleware.GetDB(r.Context())
	if pool == nil {
		return nil, handlers.NewHTTPError(http.StatusInternalServerError, "database unavailable")
	}

	repo := repository.NewCreditRepository(pool)

	credit, err := repo.GetByID(r.Context(), id)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, handlers.NewHTTPError(http.StatusNotFound, "credit not found")
		}
		return nil, handlers.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return credit, nil
}