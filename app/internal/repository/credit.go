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
	creditsHash = "credits"
	creditsList = "credits:list"
)

type CreditRepository struct {
	*baseRepo.BaseRepository
	crud *baseRepo.CRUD[domain.Credit]
}

func NewCreditRepository(db *pgxpool.Pool) *CreditRepository {
	return &CreditRepository{
		BaseRepository: baseRepo.NewBaseRepository(db),
		crud:           baseRepo.NewCRUD[domain.Credit](db, "credits"),
	}
}

func scanCredit(row pgx.Row) (domain.Credit, error) {
	var credit domain.Credit
	err := row.Scan(&credit.ID, &credit.ClientID, &credit.BankID,
		&credit.MinPayment, &credit.MaxPayment, &credit.TermMonths,
		&credit.CreditType, &credit.Status, &credit.CreatedAt)

	return credit, err
}

func (r *CreditRepository) Create(ctx context.Context, credit *domain.Credit) error {
	query := `INSERT INTO credits (client_id, bank_id, min_payment, max_payment,
							term_months, credit_type, status, created_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`

	err := r.DB().QueryRow(ctx, query, credit.ClientID, credit.BankID,
		credit.MinPayment, credit.MaxPayment, credit.TermMonths,
		credit.CreditType, credit.Status, credit.CreatedAt).Scan(&credit.ID)
    if err != nil {
        return r.HandleError(err)
    }

	// Cache
    data, _ := json.Marshal(credit)
    r.Redis().HSet(ctx, creditsHash, strconv.Itoa(credit.ID), data)
    r.Redis().ZAdd(ctx, creditsList, redis.Z{Score: float64(credit.CreatedAt.Unix()), Member: strconv.Itoa(credit.ID)})

	return r.HandleError(err)
}

func (r *CreditRepository) GetByID(ctx context.Context, id int) (*domain.Credit, error) {
	// Try cache
    data, err := r.Redis().HGet(ctx, creditsHash, strconv.Itoa(id)).Bytes()
    if err == nil {
        var credit domain.Credit
        if json.Unmarshal(data, &credit) == nil {
            return &credit, nil
        }
    }

	// DB
	credit, err := r.crud.GetByID(ctx, id, scanCredit)

	if err != nil {
		return nil, err
	}

    // Cache
	if data, err := json.Marshal(credit); err == nil {
		r.Redis().HSet(ctx, creditsHash, strconv.Itoa(id), data)
	}

	return &credit, nil
}

func (r *CreditRepository) List(ctx context.Context, pagination baseRepo.PaginationParams) (baseRepo.PaginatedResult[domain.Credit], error) {
	return r.crud.List(ctx, pagination, scanCredit, "", "created_at DESC")
}

func (r *CreditRepository) Update(ctx context.Context, credit *domain.Credit) error {
	query := `UPDATE credits
			  SET min_payment = $1, max_payment = $2, term_months = $3,
				  credit_type = $4, status = $5
			  WHERE id = $6`

	result, err := r.DB().Exec(ctx, query, credit.MinPayment, credit.MaxPayment,
		credit.TermMonths, credit.CreditType, credit.Status, credit.ID)

	if err != nil {
		return r.HandleError(err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

    data, _ := json.Marshal(credit)
	r.Redis().HSet(ctx, creditsHash, credit.ID, data)

	return nil
}

func (r *CreditRepository) Delete(ctx context.Context, id int) error {
	err := r.crud.Delete(ctx, id)
    if err != nil {
        return err
    }

    // Cache
    r.Redis().HDel(ctx, creditsHash, strconv.Itoa(id))
    r.Redis().ZRem(ctx, creditsList, strconv.Itoa(id))

    return nil
}

func (r *CreditRepository) ListByClient(ctx context.Context, clientID int, pagination baseRepo.PaginationParams) (baseRepo.PaginatedResult[domain.Credit], error) {
	whereClause := "client_id = $1"
	return r.crud.List(ctx, pagination, scanCredit, whereClause, "created_at DESC", clientID)
}

func (r *CreditRepository) ListByStatus(ctx context.Context, status string, pagination baseRepo.PaginationParams) (baseRepo.PaginatedResult[domain.Credit], error) {
	whereClause := "status = $1"
	return r.crud.List(ctx, pagination, scanCredit, whereClause, "created_at DESC", status)
}

func (r *CreditRepository) CountByStatus(ctx context.Context, status string) (int64, error) {
	return r.crud.Count(ctx, "status = $1", status)
}