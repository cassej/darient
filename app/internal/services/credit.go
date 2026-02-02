package services

import (
	"context"
	"time"

	"api/internal/domain"
	"api/internal/middleware"
	"api/internal/repository"
	baseRepo "api/pkg/repository"
)

var CreditService = creditService{}

type creditService struct{}

func (creditService) Create(ctx context.Context, clientID, bankID int, minPayment, maxPayment float64, termMonths int, creditType string) (*domain.Credit, error) {
	credit := &domain.Credit{
		ClientID:   clientID,
		BankID:     bankID,
		MinPayment: minPayment,
		MaxPayment: maxPayment,
		TermMonths: termMonths,
		CreditType: creditType,
		Status:     "PENDING",
		CreatedAt:  time.Now().UTC(),
	}

	repo := repository.NewCreditRepository(middleware.GetDB(ctx))
	return credit, repo.Create(ctx, credit)
}

func (creditService) Get(ctx context.Context, id int) (*domain.Credit, error) {
	repo := repository.NewCreditRepository(middleware.GetDB(ctx))
	return repo.GetByID(ctx, id)
}

func (creditService) Update(ctx context.Context, id int, minPayment, maxPayment *float64, termMonths *int, creditType, status *string) (*domain.Credit, error) {
	repo := repository.NewCreditRepository(middleware.GetDB(ctx))

	credit, err := repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if minPayment != nil {
		credit.MinPayment = *minPayment
	}
	if maxPayment != nil {
		credit.MaxPayment = *maxPayment
	}
	if termMonths != nil {
		credit.TermMonths = *termMonths
	}
	if creditType != nil {
		credit.CreditType = *creditType
	}
	if status != nil {
		credit.Status = *status
	}

	return credit, repo.Update(ctx, credit)
}

func (creditService) Delete(ctx context.Context, id int) error {
	repo := repository.NewCreditRepository(middleware.GetDB(ctx))
	return repo.Delete(ctx, id)
}

func (creditService) List(ctx context.Context, page, pageSize int) (interface{}, error) {
	repo := repository.NewCreditRepository(middleware.GetDB(ctx))
	pagination := baseRepo.NewPaginationParams(page, pageSize)
	return repo.List(ctx, pagination)
}