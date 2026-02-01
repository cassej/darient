package clients

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"api/internal/domain"
	"api/internal/contracts"
	"api/internal/handlers"
	"api/internal/middleware"
	"api/internal/repository"
)

func init() {
	handlers.Register("DELETE", "/clients/{id}", DeleteClient)
}

func DeleteClient(w http.ResponseWriter, r *http.Request) (interface{}, error) {
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

	repo := repository.NewClientRepository(pool)

	if err := repo.Delete(r.Context(), id); err != nil {
		if err == domain.ErrNotFound {
			return nil, handlers.NewHTTPError(http.StatusNotFound, "client not found")
		}
		return nil, handlers.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	w.WriteHeader(http.StatusNoContent)
	return nil, nil
}