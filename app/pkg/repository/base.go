package repository

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"


	"api/internal/domain"
	"api/pkg/database"
)


type BaseRepository struct {
	db *pgxpool.Pool
}

func NewBaseRepository(db *pgxpool.Pool) *BaseRepository {
	return &BaseRepository{db: db}
}

func (r *BaseRepository) Redis() *redis.Client {
	return database.Redis()
}

func (r *BaseRepository) DB() *pgxpool.Pool {
	return r.db
}

func (r *BaseRepository) HandleError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrNotFound
	}

	if pgErr, ok := err.(*pgconn.PgError); ok {
		switch pgErr.Code {
		case "23505": // unique_violation
			return domain.ErrAlreadyExists
		case "23503": // foreign_key_violation
			return domain.ErrForeignKey
		case "23502": // not_null_violation
			return domain.ErrInvalidInput
		case "23514": // check_violation
			return domain.ErrInvalidInput
		}
	}

	return err
}
