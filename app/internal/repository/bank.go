package repository

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"api/internal/domain"
	baseRepo "api/pkg/repository"
)

const (
	banksHash = "banks"
	banksList = "banks:list"
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
	query := `INSERT INTO banks (name, type, created_at) VALUES ($1, $2, $3) RETURNING id`

	err := r.DB().QueryRow(ctx, query, bank.Name, bank.Type, bank.CreatedAt).Scan(&bank.ID)
    if err != nil {
        return r.HandleError(err)
    }

	// Cache
	data, _ := json.Marshal(bank)
	r.Redis().HSet(ctx, banksHash, strconv.Itoa(bank.ID), data)
	r.Redis().ZAdd(ctx, banksList, redis.Z{Score: float64(bank.CreatedAt.Unix()), Member: strconv.Itoa(bank.ID)})

	return nil
}

func (r *BankRepository) GetByID(ctx context.Context, id int) (*domain.Bank, error) {
	// Try cache
	data, err := r.Redis().HGet(ctx, banksHash, strconv.Itoa(id)).Bytes()
	if err == nil {
		var bank domain.Bank
		if json.Unmarshal(data, &bank) == nil {
			return &bank, nil
		}
	}

	// DB
	bank, err := r.crud.GetByID(ctx, id, scanBank)

	if err != nil {
		return nil, err
	}

	// Cache
	if data, err := json.Marshal(bank); err == nil {
		r.Redis().HSet(ctx, banksHash, strconv.Itoa(id), data)
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

	// Cache
	data, _ := json.Marshal(bank)
	r.Redis().HSet(ctx, banksHash, strconv.Itoa(bank.ID), data)

	return nil
}

func (r *BankRepository) Delete(ctx context.Context, id int) error {
	err := r.crud.Delete(ctx, id)
	if err != nil {
		return err
	}

	// Cache
	r.Redis().HDel(ctx, banksHash, strconv.Itoa(id))
	r.Redis().ZRem(ctx, banksList, strconv.Itoa(id))

	return nil
}