package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ScanFunc[T any] func(row pgx.Row) (T, error)

type CRUD[T any] struct {
	*BaseRepository
	tableName string
}

func NewCRUD[T any](db *pgxpool.Pool, tableName string) *CRUD[T] {
	return &CRUD[T]{
		BaseRepository: NewBaseRepository(db),
		tableName:      tableName,
	}
}

// GetByID retrieves a single record by ID
func (c *CRUD[T]) GetByID(ctx context.Context, id string, scanFn ScanFunc[T]) (T, error) {
	var zero T
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", c.tableName)
	
	row := c.db.QueryRow(ctx, query, id)
	result, err := scanFn(row)
	
	if err != nil {
		return zero, c.HandleError(err)
	}
	
	return result, nil
}

// Delete removes a record by ID
func (c *CRUD[T]) Delete(ctx context.Context, id string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", c.tableName)
	result, err := c.db.Exec(ctx, query, id)
	
	if err != nil {
		return c.HandleError(err)
	}
	
	if result.RowsAffected() == 0 {
		return c.HandleError(pgx.ErrNoRows)
	}
	
	return nil
}

// Count returns the total number of records
func (c *CRUD[T]) Count(ctx context.Context, whereClause string, args ...any) (int64, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", c.tableName)
	if whereClause != "" {
		query += " WHERE " + whereClause
	}
	
	var count int64

	err := c.db.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, c.HandleError(err)
	}
	
	return count, nil
}

// List retrieves multiple records with pagination
func (c *CRUD[T]) List(ctx context.Context, pagination PaginationParams, scanFn ScanFunc[T], whereClause string, orderBy string, args ...any) (PaginatedResult[T], error) {
	total, err := c.Count(ctx, whereClause, args...)
	if err != nil {
		return PaginatedResult[T]{}, err
	}
	
	query := fmt.Sprintf("SELECT * FROM %s", c.tableName)
	if whereClause != "" {
		query += " WHERE " + whereClause
	}

	if orderBy != "" {
		query += " ORDER BY " + orderBy
	} else {
		query += " ORDER BY created_at DESC"
	}

	query += fmt.Sprintf(" LIMIT %d OFFSET %d", pagination.Limit(), pagination.Offset())
	
	rows, err := c.db.Query(ctx, query, args...)
	if err != nil {
		return PaginatedResult[T]{}, c.HandleError(err)
	}

	defer rows.Close()
	
	items := make([]T, 0)
	for rows.Next() {
		item, err := scanFn(rows)

		if err != nil {
			return PaginatedResult[T]{}, c.HandleError(err)
		}

		items = append(items, item)
	}
	
	if err := rows.Err(); err != nil {
		return PaginatedResult[T]{}, c.HandleError(err)
	}
	
	return NewPaginatedResult(items, total, pagination), nil
}

// Exists checks if a record with given ID exists
func (c *CRUD[T]) Exists(ctx context.Context, id string) (bool, error) {
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE id = $1)", c.tableName)
	
	var exists bool

	err := c.db.QueryRow(ctx, query, id).Scan(&exists)
	if err != nil {
		return false, c.HandleError(err)
	}
	
	return exists, nil
}
