package banks

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"

	"api/internal/handlers"
	"api/internal/contracts"
	"api/internal/contracts/banks"
	"api/internal/domain"
	"api/internal/middleware"
	"api/internal/repository"
)

func init() {
	handlers.Register("POST", "/banks", CreateBank)
}

func CreateBank(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	var input map[string]any
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return nil, handlers.NewHTTPError(http.StatusBadRequest, "invalid json format")
	}

	validated, err := contracts.Validate(input, banks.Create)
	if err != nil {
		return nil, err
	}

	name := validated["name"].(string)
	bankType := validated["type"].(string)

	bank := &domain.Bank{
		ID:        uuid.NewString(),
		Name:      name,
		Type:      bankType,
		CreatedAt: time.Now().UTC(),
	}

	pool := middleware.GetDB(r.Context())
	if pool == nil {
		return nil, handlers.NewHTTPError(http.StatusInternalServerError, "database unavailable")
	}

	repo := repository.NewBankRepository(pool)

	if err := repo.Create(r.Context(), bank); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return nil, handlers.NewHTTPError(http.StatusConflict, "bank with this name already exists")
		}
		return nil, handlers.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	w.WriteHeader(http.StatusCreated)
	return bank, nil
}