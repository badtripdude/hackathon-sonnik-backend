package postgres

import (
	"context"
	"time"

	"github.com/badtripdude/hackathon-sonnik-backend/services/robokassa/internal/models"
	"github.com/jmoiron/sqlx"
)

type PgPaymentRepo struct {
	db           *sqlx.DB
	queryTimeout time.Duration
}

func NewPgPaymentRepo(db *sqlx.DB, queryTimeout time.Duration) *PgPaymentRepo {
	return &PgPaymentRepo{db: db, queryTimeout: queryTimeout}
}

func (pg *PgPaymentRepo) Save(ctx context.Context, p *models.Payment) error {
	timeout, cancel := context.WithTimeout(ctx, pg.queryTimeout)
	defer cancel()

	rows, err := pg.db.NamedQueryContext(timeout, savePaymentQuery, p)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&p.CreatedAt, &p.UpdatedAt); err != nil {
			return err
		}
	}

	return nil
}

func (pg *PgPaymentRepo) GetByID(ctx context.Context, id string) (*models.Payment, error) {
	p := &models.Payment{}
	timeout, cancel := context.WithTimeout(ctx, pg.queryTimeout)
	defer cancel()

	if err := pg.db.GetContext(timeout, p, getPaymentByIDQuery, id); err != nil {
		return nil, err
	}

	return p, nil
}

func (pg *PgPaymentRepo) GetByUserID(ctx context.Context, userID string) ([]*models.Payment, error) {
	var payments []*models.Payment
	timeout, cancel := context.WithTimeout(ctx, pg.queryTimeout)
	defer cancel()

	if err := pg.db.SelectContext(timeout, &payments, getPaymentsByUserIDQuery, userID); err != nil {
		return nil, err
	}

	return payments, nil
}

func (pg *PgPaymentRepo) UpdateStatus(ctx context.Context, id string, status models.PaymentStatus) error {
	timeout, cancel := context.WithTimeout(ctx, pg.queryTimeout)
	defer cancel()

	_, err := pg.db.ExecContext(timeout, updatePaymentStatusQuery, status, id)
	return err
}
