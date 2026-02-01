package banks

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
	handlers.Register("DELETE", "/banks/{id}", DeleteBank)
}

func DeleteBank(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	id := chi.URLParam(r, "id")
	if id == "" {
		return nil, domain.ErrInvalidInput
	}

    err := contracts.ValidateField("id", id, contracts.FieldSpec{
        Type: "int",
        Min:  1,
    })

    if err != nil {
        return nil, handlers.NewHTTPError(http.StatusBadRequest, err.Error())
    }

	pool := middleware.GetDB(r.Context())
	if pool == nil {
		return nil, handlers.NewHTTPError(http.StatusInternalServerError, "database unavailable")
	}

	repo := repository.NewBankRepository(pool)

	if err := repo.Delete(r.Context(), id); err != nil {
		if err == domain.ErrNotFound {
			return nil, handlers.NewHTTPError(http.StatusNotFound, "bank not found")
		}
		return nil, handlers.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	w.WriteHeader(http.StatusNoContent)
	return nil, nil
}