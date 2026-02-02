package services

import (
	"context"
	"time"

	"api/internal/domain"
	"api/internal/middleware"
	"api/internal/repository"
	"api/internal/events"
	baseRepo "api/pkg/repository"
)

var CreditService = creditService{}

type creditService struct{}

func (s creditService) Create(ctx context.Context, clientID, bankID int, minPayment, maxPayment float64, termMonths int, creditType string) (*domain.Credit, error) {
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

	if err := repo.Create(ctx, credit); err != nil {
		return nil, err
	}

	if pub := middleware.GetPublisher(ctx); pub != nil {
		pub.Publish(ctx, events.Event{
			Type:      "CreditCreated",
			Timestamp: time.Now(),
			Payload: events.CreditCreatedEvent{
				CreditID:   credit.ID,
				ClientID:   credit.ClientID,
				BankID:     credit.BankID,
				CreditType: credit.CreditType,
			},
		})
	}

	return credit, nil
}

func (s creditService) Get(ctx context.Context, id int) (*domain.Credit, error) {
	repo := repository.NewCreditRepository(middleware.GetDB(ctx))
	return repo.GetByID(ctx, id)
}

func (s creditService) Update(ctx context.Context, id int, minPayment, maxPayment *float64, termMonths *int, creditType, status *string) (*domain.Credit, error) {
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

    if err := repo.Update(ctx, credit); err != nil {
		return nil, err
	}

    if credit.Status == "APPROVED" {
        if pub := middleware.GetPublisher(ctx); pub != nil {
            pub.Publish(ctx, events.Event{
                Type:      "CreditApproved",
                Timestamp: time.Now(),
                Payload: events.CreditApprovedEvent{
                    CreditID:   credit.ID,
                    ClientID:   credit.ClientID,
                    ApprovedAt: time.Now(),
                },
            })
        }
    }

	return credit, nil
}

func (s creditService) Delete(ctx context.Context, id int) error {
	repo := repository.NewCreditRepository(middleware.GetDB(ctx))
	return repo.Delete(ctx, id)
}

func (s creditService) List(ctx context.Context, page, pageSize int) (interface{}, error) {
	repo := repository.NewCreditRepository(middleware.GetDB(ctx))
	pagination := baseRepo.NewPaginationParams(page, pageSize)
	return repo.List(ctx, pagination)
}