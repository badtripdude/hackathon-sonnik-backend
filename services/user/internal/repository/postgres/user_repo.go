package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/badtripdude/hackathon-sonnik-backend/services/user/internal/models"
	errs "github.com/badtripdude/hackathon-sonnik-backend/services/user/pkg/errors"
	"github.com/badtripdude/hackathon-sonnik-backend/services/user/pkg/postgres"
)

type PgUserRepo struct {
	db           postgres.DBTX
	queryTimeout time.Duration
}

func NewPgUserRepo(db postgres.DBTX, queryTimeout time.Duration) *PgUserRepo {
	return &PgUserRepo{
		db:           db,
		queryTimeout: queryTimeout,
	}
}

func (pg *PgUserRepo) Save(ctx context.Context, user *models.User) error {
	timeout, cancel := context.WithTimeout(ctx, pg.queryTimeout)
	defer cancel()

	err := pg.db.QueryRowxContext(timeout, saveUserQuery,
		user.Email,
		user.PasswordHash,
		user.Username,
		user.BirthDate,
	).Scan(&user.ID)
	if err != nil {
		return postgres.CheckUnique(err)
	}

	return nil
}

func (pg *PgUserRepo) Update(ctx context.Context, id string, email *string, username *string) (*models.User, error) {
	timeout, cancel := context.WithTimeout(ctx, pg.queryTimeout)
	defer cancel()

	user := &models.User{}

	err := pg.db.QueryRowxContext(timeout,
		updateUserQuery,
		email,
		username,
		id,
	).StructScan(user)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}

	return user, nil
}

func (pg *PgUserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}

	timeout, cancel := context.WithTimeout(ctx, pg.queryTimeout)
	defer cancel()

	err := pg.db.GetContext(timeout, user, getUserByEmailQuery, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}

	return user, nil
}

func (pg *PgUserRepo) GetByID(ctx context.Context, id string) (*models.User, error) {
	user := &models.User{}

	timeout, cancel := context.WithTimeout(ctx, pg.queryTimeout)
	defer cancel()

	err := pg.db.GetContext(timeout, user, getUserByIDQuery, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}

	return user, nil
}

func (pg *PgUserRepo) DeleteByID(ctx context.Context, id string) error {
	timeout, cancel := context.WithTimeout(ctx, pg.queryTimeout)
	defer cancel()

	r, err := pg.db.ExecContext(timeout, deleteUserByIDQuery, id)
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

	return err
}
