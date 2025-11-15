package service

import (
	"context"
	"log"
	"time"

	"github.com/badtripdude/hackathon-sonnik-backend/services/robokassa/internal/models"
	"github.com/google/uuid"
)

type PaymentRepo interface {
	Save(ctx context.Context, p *models.Payment) error
	GetByID(ctx context.Context, id string) (*models.Payment, error)
	GetByUserID(ctx context.Context, userID string) ([]*models.Payment, error)
	UpdateStatus(ctx context.Context, id string, status models.PaymentStatus) error
}

type PaymentService struct {
	repo PaymentRepo
}

func NewPaymentService(repo PaymentRepo) *PaymentService {
	return &PaymentService{repo: repo}
}

func (s *PaymentService) CreatePayment(ctx context.Context, userID string, amount float64, currency string) (*models.Payment, error) {
	p := &models.Payment{
		ID:        uuid.New(),
		UserID:    uuid.MustParse(userID),
		Amount:    amount,
		Currency:  currency,
		Status:    models.PaymentStatusCreated,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Save(ctx, p); err != nil {
		return nil, err
	}

	// Имитация шлюза
	go s.simulatePaymentFlow(p.ID)

	return p, nil
}

func (s *PaymentService) simulatePaymentFlow(paymentID uuid.UUID) {
	// Смена статуса на pending через 5 сек
	time.Sleep(5 * time.Second)
	if err := s.repo.UpdateStatus(context.Background(), paymentID.String(), models.PaymentStatusPending); err != nil {
		log.Println("failed to set pending:", err)
		return
	}

	// Смена статуса на success через 15 сек
	time.Sleep(15 * time.Second)
	if err := s.repo.UpdateStatus(context.Background(), paymentID.String(), models.PaymentStatusSuccess); err != nil {
		log.Println("failed to set success:", err)
		return
	}

	// Таймаут на завершение платежа (если не успели)
	time.AfterFunc(60*time.Second, func() {
		p, err := s.repo.GetByID(context.Background(), paymentID.String())
		if err != nil {
			return
		}
		if p.Status != models.PaymentStatusSuccess && p.Status != models.PaymentStatusFailed {
			_ = s.repo.UpdateStatus(context.Background(), paymentID.String(), models.PaymentStatusFailed)
		}
	})
}

func (s *PaymentService) GetPaymentByID(ctx context.Context, id string) (*models.Payment, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *PaymentService) GetPaymentsByUser(ctx context.Context, userID string) ([]*models.Payment, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *PaymentService) UpdatePaymentStatus(ctx context.Context, id string, status models.PaymentStatus) error {
	return s.repo.UpdateStatus(ctx, id, status)
}
