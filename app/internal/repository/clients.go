package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"api/internal/domain"
	baseRepo "api/pkg/repository"
)

type ClientRepository struct {
	*baseRepo.BaseRepository
	crud *baseRepo.CRUD[domain.Client]
}

func NewClientRepository(db *pgxpool.Pool) *ClientRepository {
	return &ClientRepository{
		BaseRepository: baseRepo.NewBaseRepository(db),
		crud:           baseRepo.NewCRUD[domain.Client](db, "clients"),
	}
}

func scanClient(row pgx.Row) (domain.Client, error) {
	var client domain.Client

	err := row.Scan(&client.ID, &client.FullName, &client.Email,
		&client.BirthDate, &client.Country, &client.CreatedAt)

	return client, err
}

func (r *ClientRepository) Create(ctx context.Context, client *domain.Client) error {
	query := `INSERT INTO clients (id, full_name, email, birth_date, country, created_at)
			  VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.DB().Exec(ctx, query, client.ID, client.FullName, client.Email,
		client.BirthDate, client.Country, client.CreatedAt)

	return r.HandleError(err)
}

func (r *ClientRepository) GetByID(ctx context.Context, id string) (*domain.Client, error) {
	client, err := r.crud.GetByID(ctx, id, scanClient)

	if err != nil {
		return nil, err
	}

	return &client, nil
}

func (r *ClientRepository) List(ctx context.Context, pagination baseRepo.PaginationParams) (baseRepo.PaginatedResult[domain.Client], error) {
	return r.crud.List(ctx, pagination, scanClient, "", "created_at DESC")
}

func (r *ClientRepository) Update(ctx context.Context, client *domain.Client) error {
	query := `UPDATE clients
			  SET full_name = $1, email = $2, birth_date = $3, country = $4
			  WHERE id = $5`

	result, err := r.DB().Exec(ctx, query, client.FullName, client.Email,
		client.BirthDate, client.Country, client.ID)
	if err != nil {
		return r.HandleError(err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *ClientRepository) Delete(ctx context.Context, id string) error {
	return r.crud.Delete(ctx, id)
}

