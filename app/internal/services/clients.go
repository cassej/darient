package services

import (
	"context"
	"time"

	"api/internal/domain"
	"api/internal/middleware"
	"api/internal/repository"
	baseRepo "api/pkg/repository"
)

var ClientService = clientService{}

type clientService struct{}

func (clientService) Create(ctx context.Context, fullName, email, birthDate, country string) (*domain.Client, error) {
	client := &domain.Client{
		FullName:  fullName,
		Email:     email,
		BirthDate: birthDate,
		Country:   country,
		CreatedAt: time.Now().UTC(),
	}

	repo := repository.NewClientRepository(middleware.GetDB(ctx))
	return client, repo.Create(ctx, client)
}

func (clientService) Get(ctx context.Context, id int) (*domain.Client, error) {
	repo := repository.NewClientRepository(middleware.GetDB(ctx))
	return repo.GetByID(ctx, id)
}

func (clientService) Update(ctx context.Context, id int, fullName, email, birthDate, country *string) (*domain.Client, error) {
	repo := repository.NewClientRepository(middleware.GetDB(ctx))

	client, err := repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if fullName != nil {
		client.FullName = *fullName
	}
	if email != nil {
		client.Email = *email
	}
	if birthDate != nil {
		client.BirthDate = *birthDate
	}
	if country != nil {
		client.Country = *country
	}

	return client, repo.Update(ctx, client)
}

func (clientService) Delete(ctx context.Context, id int) error {
	repo := repository.NewClientRepository(middleware.GetDB(ctx))
	return repo.Delete(ctx, id)
}

func (clientService) List(ctx context.Context, page, pageSize int) (interface{}, error) {
	repo := repository.NewClientRepository(middleware.GetDB(ctx))
	pagination := baseRepo.NewPaginationParams(page, pageSize)
	return repo.List(ctx, pagination)
}