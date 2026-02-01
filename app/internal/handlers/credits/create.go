package credits

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"

	"api/internal/handlers"
	"api/internal/contracts"
	"api/internal/contracts/credits"
	"api/internal/domain"
	"api/internal/middleware"
	"api/internal/repository"
)

func init() {
	handlers.Register("POST", "/credits", CreateCredit)
}

func CreateCredit(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	var input map[string]any
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return nil, handlers.NewHTTPError(http.StatusBadRequest, "invalid json format")
	}

	validated, err := contracts.Validate(input, credits.Create)
	if err != nil {
		return nil, err
	}

	credit := &domain.Credit{
		ID:        uuid.NewString(),
		ClientID:   validated["client_id"].(string),
		BankID:     validated["bank_id"].(string),
		MinPayment: validated["min_payment"].(float64),
		MaxPayment: validated["max_payment"].(float64),
		TermMonths: validated["term_months"].(int),
		CreditType: validated["credit_type"].(string),
		Status:     "PENDING",
		CreatedAt: time.Now().UTC(),
	}

	pool := middleware.GetDB(r.Context())
	if pool == nil {
		return nil, handlers.NewHTTPError(http.StatusInternalServerError, "database unavailable")
	}

	repo := repository.NewCreditRepository(pool)

	if err := repo.Create(r.Context(), credit); err != nil {
		return nil, handlers.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	w.WriteHeader(http.StatusCreated)
	return credit, nil
}