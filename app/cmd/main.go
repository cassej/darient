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
	"api/internal/events"
	"api/internal/handlers"
	_ "api/internal/handlers/banks"
	_ "api/internal/handlers/clients"
	_ "api/internal/handlers/credits"
	mw "api/internal/middleware"
	"api/pkg/database"
)

func main() {
	cfg := config.MustLoad()

	var level slog.Level
    switch cfg.LogLevelString() {
    case "debug":
        level = slog.LevelDebug
    case "warn":
        level = slog.LevelWarn
    case "error":
        level = slog.LevelError
    default:
        level = slog.LevelInfo
    }

    log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
        Level:     level,
        AddSource: true,
    }))

    slog.SetDefault(log)

    log.Info("logger initialized", "level", level.String())

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	if err := database.ConnectRedis(ctx, cfg.GetRedisAddr()); err != nil {
		log.Error("failed to connect to redis", "err", err)
		os.Exit(1)
	}

	defer database.CloseRedis()

	go database.StartRedisHealthCheck(ctx, cfg.RedisHealthCheckInterval)

	publisher := events.NewRedisPublisher(database.Redis())

	db, err := database.Connect(ctx, cfg.GetDBDSN(), cfg.DBMaxConns, cfg.DBMinConns)

	if err != nil {
		log.Error("failed to connect to database", "err", err)
		os.Exit(1)
	}

	defer db.Close()

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(cfg.ReadHeaderTimeout))
	r.Use(mw.DBMiddleware(db))
	r.Use(mw.PublisherMiddleware(publisher))
	r.Use(mw.LoggerMiddleware(log))

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