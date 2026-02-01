package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"api/internal/domain"
	baseRepo "api/pkg/repository"
)

type BankRepository struct {
	*baseRepo.BaseRepository
	crud *baseRepo.CRUD[domain.Bank]
}

func NewBankRepository(db *pgxpool.Pool) *BankRepository {
	return &BankRepository{
		BaseRepository: baseRepo.NewBaseRepository(db),
		crud:           baseRepo.NewCRUD[domain.Bank](db, "banks"),
	}
}

func scanBank(row pgx.Row) (domain.Bank, error) {
	var bank domain.Bank
	err := row.Scan(&bank.ID, &bank.Name, &bank.Type, &bank.CreatedAt)

	return bank, err
}

func (r *BankRepository) Create(ctx context.Context, bank *domain.Bank) error {
	query := `INSERT INTO banks (id, name, type, created_at) VALUES ($1, $2, $3, $4)`

	_, err := r.DB().Exec(ctx, query, bank.ID, bank.Name, bank.Type, bank.CreatedAt)

	return r.HandleError(err)
}

func (r *BankRepository) GetByID(ctx context.Context, id string) (*domain.Bank, error) {
	bank, err := r.crud.GetByID(ctx, id, scanBank)

	if err != nil {
		return nil, err
	}

	return &bank, nil
}

func (r *BankRepository) List(ctx context.Context, pagination baseRepo.PaginationParams) (baseRepo.PaginatedResult[domain.Bank], error) {
	return r.crud.List(ctx, pagination, scanBank, "", "created_at DESC")
}

func (r *BankRepository) Update(ctx context.Context, bank *domain.Bank) error {
	query := `UPDATE banks SET name = $1, type = $2 WHERE id = $3`

	result, err := r.DB().Exec(ctx, query, bank.Name, bank.Type, bank.ID)
	if err != nil {
		return r.HandleError(err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *BankRepository) Delete(ctx context.Context, id string) error {
	return r.crud.Delete(ctx, id)
}