package clients

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"

	"api/internal/handlers"
	"api/internal/contracts"
	"api/internal/contracts/clients"
	"api/internal/domain"
	"api/internal/middleware"
	"api/internal/repository"
)

func init() {
	handlers.Register("POST", "/clients", CreateClient)
}

func CreateClient(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	var input map[string]any
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return nil, handlers.NewHTTPError(http.StatusBadRequest, "invalid json format")
	}

	validated, err := contracts.Validate(input, clients.Create)
	if err != nil {
		return nil, err
	}

	client := &domain.Client{
		ID:        uuid.NewString(),
		FullName:  validated["full_name"].(string),
		Email:     validated["email"].(string),
		BirthDate: validated["birth_date"].(string),
		Country:   validated["country"].(string),
		CreatedAt: time.Now().UTC(),
	}

	pool := middleware.GetDB(r.Context())
	if pool == nil {
		return nil, handlers.NewHTTPError(http.StatusInternalServerError, "database unavailable")
	}

	repo := repository.NewClientRepository(pool)

	if err := repo.Create(r.Context(), client); err != nil {
		return nil, handlers.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	w.WriteHeader(http.StatusCreated)
	return client, nil
}