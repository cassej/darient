package middleware

import (
	"context"
	"net/http"

	"api/internal/events"
)

const publisherKey contextKey = "publisher"

func PublisherMiddleware(pub events.EventPublisher) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), publisherKey, pub)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetPublisher(ctx context.Context) events.EventPublisher {
	if v := ctx.Value(publisherKey); v != nil {
		if pub, ok := v.(events.EventPublisher); ok {
			return pub
		}
	}
	return nil
}
