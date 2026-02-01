package clients

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"api/internal/contracts"
	"api/internal/contracts/clients"
	"api/internal/domain"
	"api/internal/handlers"
	"api/internal/middleware"
	"api/internal/repository"
)

func init() {
	handlers.Register("PUT", "/clients/{id}", UpdateClient)
}

func UpdateClient(w http.ResponseWriter, r *http.Request) (interface{}, error) {
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

	var input map[string]any
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return  nil, handlers.NewHTTPError(http.StatusBadRequest, "invalid json format")
	}

	validated, err := contracts.Validate(input, clients.Update)
	if err != nil {
		return nil, err
	}

	pool := middleware.GetDB(r.Context())
	if pool == nil {
		return nil, handlers.NewHTTPError(http.StatusInternalServerError, "database unavailable")
	}

	repo := repository.NewClientRepository(pool)

	existingClient, err := repo.GetByID(r.Context(), id)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, handlers.NewHTTPError(http.StatusNotFound, "client not found")
		}
		return nil, handlers.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if fullName, ok := validated["full_name"].(string); ok {
		existingClient.FullName = fullName
	}
	if email, ok := validated["email"].(string); ok {
		existingClient.Email = email
	}
	if birthDate, ok := validated["birth_date"].(string); ok {
		existingClient.BirthDate = birthDate
	}
	if country, ok := validated["country"].(string); ok {
		existingClient.Country = country
	}

	if err := repo.Update(r.Context(), existingClient); err != nil {
		if err == domain.ErrNotFound {
			return nil, handlers.NewHTTPError(http.StatusNotFound, "client not found")
		}
		return nil, handlers.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return existingClient, nil
}