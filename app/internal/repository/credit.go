package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"api/internal/domain"
	baseRepo "api/pkg/repository"
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
	query := `INSERT INTO credits (id, client_id, bank_id, min_payment, max_payment,
							term_months, credit_type, status, created_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.DB().Exec(ctx, query, credit.ID, credit.ClientID, credit.BankID,
		credit.MinPayment, credit.MaxPayment, credit.TermMonths,
		credit.CreditType, credit.Status, credit.CreatedAt)

	return r.HandleError(err)
}

func (r *CreditRepository) GetByID(ctx context.Context, id string) (*domain.Credit, error) {
	credit, err := r.crud.GetByID(ctx, id, scanCredit)

	if err != nil {
		return nil, err
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

	return nil
}

func (r *CreditRepository) Delete(ctx context.Context, id string) error {
	return r.crud.Delete(ctx, id)
}

func (r *CreditRepository) ListByClient(ctx context.Context, clientID string, pagination baseRepo.PaginationParams) (baseRepo.PaginatedResult[domain.Credit], error) {
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