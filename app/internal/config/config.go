package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port string

	LogLevel slog.Level

	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	ShutdownGrace     time.Duration

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	DBMaxConns int32
	DBMinConns int32
}

func (c Config) LogLevelString() string {
	switch c.LogLevel {
	case slog.LevelDebug:
		return "debug"
	case slog.LevelInfo:
		return "info"
	case slog.LevelWarn:
		return "warn"
	case slog.LevelError:
		return "error"
	default:
		return "info"
	}
}

func (c *Config) GetDBDSN() string {
    return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
        c.DBUser,
        c.DBPassword,
        c.DBHost,
        c.DBPort,
        c.DBName,
        c.DBSSLMode,
    )
}

func MustLoad() Config {
	var cfg Config

	cfg.Port = envOr("PORT", "8080")

	logLevelStr := envOr("LOG_LEVEL", "INFO")
	switch logLevelStr {
	case "debug":
		cfg.LogLevel = slog.LevelDebug
	case "info":
		cfg.LogLevel = slog.LevelInfo
	case "warn":
		cfg.LogLevel = slog.LevelWarn
	case "error":
		cfg.LogLevel = slog.LevelError
	default:
		cfg.LogLevel = slog.LevelInfo
	}

	cfg.ReadTimeout = durationEnvOr("READ_TIMEOUT_SEC", 10*time.Second)
	cfg.WriteTimeout = durationEnvOr("WRITE_TIMEOUT_SEC", 30*time.Second)
	cfg.IdleTimeout = durationEnvOr("IDLE_TIMEOUT_SEC", 60*time.Second)
	cfg.ReadHeaderTimeout = durationEnvOr("READ_HEADER_TIMEOUT_SEC", 5*time.Second)
	cfg.ShutdownGrace = durationEnvOr("SHUTDOWN_GRACE_SEC", 15*time.Second)

	cfg.DBHost = envOr("DB_HOST", "odyssey")
	cfg.DBPort = envOr("DB_PORT", "6432")
	cfg.DBUser = envOr("DB_USER", "postgres")
	cfg.DBPassword = envOr("DB_PASSWORD", "postgres")
	cfg.DBName = envOr("DB_NAME", "credits")
	cfg.DBSSLMode = envOr("DB_SSLMODE", "disable")
	cfg.DBMaxConns = int32(intEnvOr("DB_MAX_CONNS", 25))
	cfg.DBMinConns = int32(intEnvOr("DB_MIN_CONNS", 5))

	return cfg
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func intEnvOr(key string, fallback int) int {
	s := envOr(key, "")
	if s == "" {
		return fallback
	}
	val, err := strconv.Atoi(s)
	if err != nil {
		return fallback
	}
	return val
}

func durationEnvOr(key string, fallback time.Duration) time.Duration {
	s := envOr(key, "")
	if s == "" {
		return fallback
	}
	sec, err := strconv.Atoi(s)
	if err != nil {
		return fallback
	}
	return time.Duration(sec) * time.Second
}