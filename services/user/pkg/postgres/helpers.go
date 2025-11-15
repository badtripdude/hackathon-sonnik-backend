package postgres

import (
	"errors"
	"fmt"
	errs "github.com/badtripdude/hackathon-sonnik-backend/services/user/pkg/errors"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	CodeUniqueViolation = "23505"
)

func CheckUnique(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == CodeUniqueViolation {
		return fmt.Errorf("%w: %s must be unique", errs.ErrAlreadyExists, pgErr.ConstraintName)
	}

	return err
}
