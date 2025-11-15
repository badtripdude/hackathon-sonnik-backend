package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/badtripdude/hackathon-sonnik-backend/services/user/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func NewPgClient(ctx context.Context, cfg *config.Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx", cfg.PgDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	db.SetMaxOpenConns(cfg.PgMaxOpenConns)
	db.SetMaxIdleConns(cfg.PgMaxIdleConns)
	db.SetConnMaxLifetime(cfg.PgConnMaxLifetime)

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := db.PingContext(pingCtx); err != nil {
		return nil, fmt.Errorf("postgres ping failed: %w", err)
	}

	return db, nil
}
