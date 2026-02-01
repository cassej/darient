package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"api/internal/domain"
)

type BankRepository struct {
	db *pgxpool.Pool
}

func NewBankRepository(db *pgxpool.Pool) *BankRepository {
	return &BankRepository{db: db}
}

func (r *BankRepository) Create(ctx context.Context, bank *domain.Bank) error {
	query := `INSERT INTO banks (id, name, type, created_at) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(ctx, query, bank.ID, bank.Name, bank.Type, bank.CreatedAt)
	return err
}

func (r *BankRepository) GetByID(ctx context.Context, id string) (*domain.Bank, error) {
	query := `SELECT id, name, type, created_at FROM banks WHERE id = $1`

	var bank domain.Bank
	err := r.db.QueryRow(ctx, query, id).Scan(&bank.ID, &bank.Name, &bank.Type, &bank.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return &bank, nil
}

func (r *BankRepository) List(ctx context.Context) ([]domain.Bank, error) {
	query := `SELECT id, name, type, created_at FROM banks ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var banks []domain.Bank
	for rows.Next() {
		var bank domain.Bank

		if err := rows.Scan(&bank.ID, &bank.Name, &bank.Type, &bank.CreatedAt); err != nil {
			return nil, err
		}

		banks = append(banks, bank)
	}

	return banks, rows.Err()
}

func (r *BankRepository) Update(ctx context.Context, bank *domain.Bank) error {
	query := `UPDATE banks SET name = $1, type = $2 WHERE id = $3`
	result, err := r.db.Exec(ctx, query, bank.Name, bank.Type, bank.ID)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *BankRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM banks WHERE id = $1`
	result, err := r.db.Exec(ctx, query, id)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	return nil
}