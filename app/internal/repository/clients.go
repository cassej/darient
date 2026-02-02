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
	clientsHash = "clients"
	clientsList = "clients:list"
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
	query := `INSERT INTO clients (full_name, email, birth_date, country, created_at)
			  VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err := r.DB().QueryRow(ctx, query, client.FullName, client.Email,
		client.BirthDate, client.Country, client.CreatedAt).Scan(&client.ID)
	if err != nil {
        return r.HandleError(err)
    }

	// Cache
	data, _ := json.Marshal(client)
	r.Redis().HSet(ctx, clientsHash, strconv.Itoa(client.ID), data)
	r.Redis().ZAdd(ctx, clientsList, redis.Z{Score: float64(client.CreatedAt.Unix()), Member: strconv.Itoa(client.ID)})

	return nil
}

func (r *ClientRepository) GetByID(ctx context.Context, id int) (*domain.Client, error) {
	// Try cache
	data, err := r.Redis().HGet(ctx, clientsHash, strconv.Itoa(id)).Bytes()
	if err == nil {
		var client domain.Client
		if json.Unmarshal(data, &client) == nil {
			return &client, nil
		}
	}

	// DB
	client, err := r.crud.GetByID(ctx, id, scanClient)

	if err != nil {
		return nil, err
	}

	// Cache
	if data, err := json.Marshal(client); err == nil {
		r.Redis().HSet(ctx, clientsHash, strconv.Itoa(id), data)
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

	// Cache
	data, _ := json.Marshal(client)
	r.Redis().HSet(ctx, clientsHash, strconv.Itoa(client.ID), data)

	return nil
}

func (r *ClientRepository) Delete(ctx context.Context, id int) error {
	err := r.crud.Delete(ctx, id)
	if err != nil {
		return err
    }

	// Cache
	r.Redis().HDel(ctx, clientsHash, strconv.Itoa(id))
	r.Redis().ZRem(ctx, clientsList, strconv.Itoa(id))

	return nil
}