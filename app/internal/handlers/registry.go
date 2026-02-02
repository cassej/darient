package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"github.com/go-chi/chi/v5"

	"api/internal/contracts"
	"api/internal/domain"
)

type HandlerFunc func(ctx context.Context, data map[string]any) (interface{}, error)

type Route struct {
	Method  string
	Path    string
	Handler http.HandlerFunc
}

var routes []Route

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
            case "PATCH":
                r.Patch(route.Path, route.Handler)
		}
	}
}

func Register(contract contracts.Contract, handler HandlerFunc) {
	routes = append(routes, Route{
		Method:  contract.Method,
		Path:    contract.URI,
		Handler: wrapWithValidation(handler, contract),
	})
}

func wrapWithValidation(fn HandlerFunc, contract contracts.Contract) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		validated := make(map[string]any)

		// Validate URI params
		for _, param := range contract.URIParams() {
			value := chi.URLParam(r, param)
			if value == "" {
				writeError(w, http.StatusBadRequest, "missing "+param+" parameter")
				return
			}

			spec, ok := contract.Required[param]
			if !ok {
				spec, ok = contract.Optional[param]
			}
			if !ok {
				writeError(w, http.StatusBadRequest, "unknown parameter: "+param)
				return
			}

			if err := contracts.ValidateField(param, value, spec); err != nil {
				writeError(w, http.StatusBadRequest, err.Error())
				return
			}

			validated[param] = contracts.Normalize(value, spec)
		}

		if r.ContentLength > 0 {
			var input map[string]any
			if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
				writeError(w, http.StatusBadRequest, "invalid json format")
				return
			}

			bodyValidated, err := contracts.Validate(input, contract)
			if err != nil {
				writeError(w, http.StatusBadRequest, err.Error())
				return
			}

			for k, v := range bodyValidated {
				validated[k] = v
			}
		}

        w.Header().Set("Content-Type", "application/json")

		data, err := fn(ctx, validated)

		if err != nil {
			handleError(w, err)
			return
		}

		// Determine status code
		statusCode := http.StatusOK
		if contract.Method == "POST" && data != nil {
			statusCode = http.StatusCreated
		}
		if data == nil {
			statusCode = http.StatusNoContent
		}

		w.WriteHeader(statusCode)
		if data != nil {
			json.NewEncoder(w).Encode(data)
		}
	}
}

func handleError(w http.ResponseWriter, err error) {
    slog.Error("handler error", "err", err)

    if httpErr, ok := err.(*HTTPError); ok {
        writeError(w, httpErr.Status, httpErr.Message)
        return
    }

    if err == domain.ErrNotFound {
        writeError(w, http.StatusNotFound, "not found")
        return
    }

    if err == domain.ErrInvalidInput {
        writeError(w, http.StatusBadRequest, "invalid input")
        return
    }

	if err == domain.ErrAlreadyExists {
		writeError(w, http.StatusConflict, "already exists")
			return
		}

	writeError(w, http.StatusInternalServerError, err.Error())
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}