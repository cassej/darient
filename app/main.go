package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"api/internal/config"
	"api/internal/handlers"
	loggerMw "api/internal/middleware"
    "api/pkg"
)

func main() {
	cfg := config.MustLoad()

	log := logger.New(cfg.LogLevelString())
	slog.SetDefault(log)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(cfg.ReadHeaderTimeout))
	r.Use(loggerMw.Middleware(log))

	r.Group(func(r chi.Router) {
        r.Use(func(next http.Handler) http.Handler {
            return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                if os.Getenv("PPROF_ENABLED") != "true" {
                    http.Error(w, "pprof disabled", http.StatusForbidden)
                    return
                }
                next.ServeHTTP(w, r)
            })
        })
		r.HandleFunc("/debug/pprof/*", http.DefaultServeMux.ServeHTTP)
    })

    handlers.RegisterAll(r)

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           r,
		ReadTimeout:       cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
	}

	go func() {
		log.Info("server listening", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("server failed", "err", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	log.Info("shutdown signal received")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.ShutdownGrace)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("graceful shutdown failed", "err", err)
	}

	log.Info("server stopped")
}