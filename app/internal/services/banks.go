package services

import (
	"context"
	"time"

	"api/internal/domain"
	"api/internal/middleware"
	"api/internal/repository"
)

var BankService = bankService{}

type bankService struct{}

func (bankService) Create(ctx context.Context, name, bankType string) (*domain.Bank, error) {
	bank := &domain.Bank{
		Name:      name,
		Type:      bankType,
		CreatedAt: time.Now().UTC(),
	}

	pool := middleware.GetDB(ctx)
	repo := repository.NewBankRepository(pool)

	if err := repo.Create(ctx, bank); err != nil {
		return nil, err
	}

	return bank, nil
}

func (bankService) Get(ctx context.Context, id string) (*domain.Bank, error) {
	pool := middleware.GetDB(ctx)
	repo := repository.NewBankRepository(pool)
	return repo.GetByID(ctx, id)
}

func (bankService) Update(ctx context.Context, id string, name, bankType *string) (*domain.Bank, error) {
	pool := middleware.GetDB(ctx)
	repo := repository.NewBankRepository(pool)

	bank, err := repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if name != nil {
		bank.Name = *name
	}
	if bankType != nil {
		bank.Type = *bankType
	}

	if err := repo.Update(ctx, bank); err != nil {
		return nil, err
	}

	return bank, nil
}

func (bankService) Delete(ctx context.Context, id string) error {
	pool := middleware.GetDB(ctx)
	repo := repository.NewBankRepository(pool)
	return repo.Delete(ctx, id)
}

func (bankService) List(ctx context.Context, page, pageSize int) (interface{}, error) {
	pool := middleware.GetDB(ctx)
	repo := repository.NewBankRepository(pool)

	pagination := repository.NewPaginationParams(page, pageSize)
	return repo.List(ctx, pagination)
}