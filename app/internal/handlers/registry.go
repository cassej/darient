package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"github.com/go-chi/chi/v5"
)

type Route struct {
	Method  string
	Path    string
	Handler http.HandlerFunc
}

type HandlerFunc func(w http.ResponseWriter, r *http.Request) (interface{}, error)

var routes []Route

func Register(method, path string, handler HandlerFunc) {
	routes = append(routes, Route{
		Method: method,
		Path:   path,
		Handler: wrap(handler),
	})
}

func wrap(fn HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if w.Header().Get("Content-Type") == "" {
			w.Header().Set("Content-Type", "application/json")
		}

		data, err := fn(w, r)
		if err != nil {
			slog.Error("handler error", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		if data != nil {
			json.NewEncoder(w).Encode(data)
		}
	}
}

func RegisterAll(r chi.Router) {
	for _, route := range routes {
		switch route.Method {
		case "GET":
			r.Get(route.Path, route.Handler)
		case "POST":
			r.Post(route.Path, route.Handler)
		case "PUT":
			r.Put(route.Path, route.Handler)
		case "DELETE":
			r.Delete(route.Path, route.Handler)
		}
	}
}