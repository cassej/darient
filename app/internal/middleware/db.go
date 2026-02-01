package middleware

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type contextKey string

const dbKey contextKey = "db"

func DBMiddleware(pool *pgxpool.Pool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), dbKey, pool)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetDB(ctx context.Context) *pgxpool.Pool {
	if v := ctx.Value(dbKey); v != nil {
		if pool, ok := v.(*pgxpool.Pool); ok {
			return pool
		}
	}

	return nil
}