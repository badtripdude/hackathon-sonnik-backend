package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/badtripdude/hackathon-sonnik-backend/services/user/internal/models"
	"github.com/badtripdude/hackathon-sonnik-backend/services/user/pkg/postgres"
	"time"

	errs "github.com/badtripdude/hackathon-sonnik-backend/services/user/pkg/errors"
)

type PgRefreshTokenRepo struct {
	db           postgres.DBTX
	queryTimeout time.Duration
}

func NewPgRefreshTokenRepo(db postgres.DBTX, queryTimeout time.Duration) *PgRefreshTokenRepo {
	return &PgRefreshTokenRepo{
		db:           db,
		queryTimeout: queryTimeout,
	}
}

func (pg *PgRefreshTokenRepo) Save(ctx context.Context, token *models.RefreshToken) error {
	timeout, cancel := context.WithTimeout(ctx, pg.queryTimeout)
	defer cancel()

	_, err := pg.db.NamedExecContext(timeout, saveRefreshTokenQuery, token)
	if err != nil {
		return err
	}

	return nil
}

func (pg *PgRefreshTokenRepo) GetByTokenHash(ctx context.Context, token string) (*models.RefreshToken, error) {
	refreshToken := &models.RefreshToken{}

	timeout, cancel := context.WithTimeout(ctx, pg.queryTimeout)
	defer cancel()

	err := pg.db.GetContext(timeout, refreshToken, getRefreshTokenByTokenHashQuery, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}

	return refreshToken, nil
}

func (pg *PgRefreshTokenRepo) RevokeByTokenHash(ctx context.Context, token string) error {
	timeout, cancel := context.WithTimeout(ctx, pg.queryTimeout)
	defer cancel()

	r, err := pg.db.ExecContext(timeout, revokeRefreshTokenByTokenHashQuery, time.Now().UTC(), token)
	if err != nil {
		return err
	}

	rows, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errs.ErrNotFound
	}

	return nil
}

func (pg *PgRefreshTokenRepo) RevokeAllByUserId(ctx context.Context, userId string) error {
	timeout, cancel := context.WithTimeout(ctx, pg.queryTimeout)
	defer cancel()

	r, err := pg.db.ExecContext(timeout, revokeAllRefreshTokenByUserIdQuery, time.Now().UTC(), userId)
	if err != nil {
		return err
	}

	rows, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errs.ErrNotFound
	}

	return nil
}
