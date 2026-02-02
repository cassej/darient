package services

import (
	"context"
	"time"
	"sync"
	"errors"

	"api/internal/domain"
	"api/internal/middleware"
	"api/internal/repository"
	"api/internal/events"
	baseRepo "api/pkg/repository"
)

var CreditService = creditService{}

type creditService struct{}

func (s creditService) Create(ctx context.Context, clientID, bankID int, minPayment, maxPayment float64, termMonths int, creditType string) (*domain.Credit, error) {
    score, err := s.ValidateEligibility(ctx, clientID, bankID)
    if err != nil {
        return nil, err
    }

    if score < 50 {
        return nil, domain.ErrNotEligible
    }

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

func (s creditService) ValidateEligibility(ctx context.Context, clientID, bankID int) (int, error) {
    ctx, cancel := context.WithCancel(ctx)
    defer cancel()

    type res struct{
        score int;
        err error
    }

    ch := make(chan res, 3)
    var wg sync.WaitGroup

    db := middleware.GetDB(ctx)

    // Check that age is between 18 and 70
    wg.Add(1)
    go func() {
        defer wg.Done()

        client, err := repository.NewClientRepository(db).GetByID(ctx, clientID)
        if err != nil {
            select {
                case ch <- res{0, err}:
                case <-ctx.Done():
            }
            return
        }

        birth, err := time.Parse("2000-01-02", client.BirthDate)
        if err != nil {
            ch <- res{0, err}
            return
        }

        age := int(time.Since(birth).Hours() / 8760)
        score := 0

        if age >= 18 {
            if age <= 70 {
                score = 35
            } else {
                score = 15
            }
        }

        select {
            case ch <- res{score, nil}:
            case <-ctx.Done():
        }
    }()

    // Check type of bank
    wg.Add(1)
    go func() {
        defer wg.Done()

        bank, err := repository.NewBankRepository(db).GetByID(ctx, bankID)
        if err != nil {
            select {
                case ch <- res{0, err}:
                case <-ctx.Done():
            }
            return
        }

        score := 0

        switch bank.Type {
            case "PRIVATE":
                score = 30
            case "GOVERNMENT":
                score = 20
        }

        select {
            case ch <- res{score, nil}:
            case <-ctx.Done():
        }
    }()

    // Check country
    wg.Add(1)
    go func() {
        defer wg.Done()

        client, err := repository.NewClientRepository(db).GetByID(ctx, clientID)
        if err != nil {
            select {
                case ch <- res{0, err}:
                case <-ctx.Done():
            }
            return
        }

        score := 10
        switch client.Country {
            case "USA", "Canada", "Chili":
                score = 35
            case "Mexico", "Brazil", "Panama":
                score = 20
        }

        select {
            case ch <- res{score, nil}:
            case <-ctx.Done():
        }
    }()

    go func() {
        wg.Wait()
        close(ch)
    }()

    var total int
    var errs []error

    done := false
    for !done {
        select {
        case r, ok := <-ch:
            if !ok {
                done = true
                break
            }
            if r.err != nil {
                errs = append(errs, r.err)
            } else {
                total += r.score
            }
        case <-ctx.Done():
            return total, ctx.Err()
        }
    }

    if len(errs) > 0 {
        return 0, errors.Join(errs...)
    }

    return total, nil
}