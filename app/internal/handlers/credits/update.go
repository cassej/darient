package credits

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"api/internal/contracts"
	"api/internal/contracts/credits"
	"api/internal/domain"
	"api/internal/handlers"
	"api/internal/middleware"
	"api/internal/repository"
)

func init() {
	handlers.Register("PUT", "/credits/{id}", UpdateCredit)
}

func UpdateCredit(w http.ResponseWriter, r *http.Request) (interface{}, error) {
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
		return nil, handlers.NewHTTPError(http.StatusBadRequest, "invalid json format")
	}

	validated, err := contracts.Validate(input, credits.Update)
	if err != nil {
		return nil, err
	}

	pool := middleware.GetDB(r.Context())
	if pool == nil {
		return nil, handlers.NewHTTPError(http.StatusInternalServerError, "database unavailable")
	}

	repo := repository.NewCreditRepository(pool)

	existingCredit, err := repo.GetByID(r.Context(), id)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, handlers.NewHTTPError(http.StatusNotFound, "credit not found")
		}
		return nil, handlers.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if minPayment, ok := validated["min_payment"].(float64); ok {
		existingCredit.MinPayment = minPayment
	}
	if maxPayment, ok := validated["max_payment"].(float64); ok {
		existingCredit.MaxPayment = maxPayment
	}
	if termMonths, ok := validated["term_months"].(int); ok {
		existingCredit.TermMonths = termMonths
	}
	if creditType, ok := validated["credit_type"].(string); ok {
		existingCredit.CreditType = creditType
	}
	if status, ok := validated["status"].(string); ok {
		existingCredit.Status = status
	}

	if err := repo.Update(r.Context(), existingCredit); err != nil {
		if err == domain.ErrNotFound {
			return nil, handlers.NewHTTPError(http.StatusNotFound, "credit not found")
		}
		return nil, handlers.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return existingCredit, nil
}