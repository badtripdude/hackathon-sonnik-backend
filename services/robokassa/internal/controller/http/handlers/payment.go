package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/badtripdude/hackathon-sonnik-backend/services/robokassa/internal/models"
	"github.com/badtripdude/hackathon-sonnik-backend/services/robokassa/internal/service"
	"github.com/gorilla/mux"
)

type PaymentHandler struct {
	svc *service.PaymentService
}

func NewPaymentHandler(svc *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{svc: svc}
}

// CreatePayment godoc
// @Summary Create a new payment
// @Description Create a payment for a user with specified amount and currency
// @Tags payments
// @Accept json
// @Produce json
// @Param request body models.CreatePaymentRequest true "Create payment request"
// @Success 200 {object} models.PaymentResponse "Payment created successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid request body"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /payments [post]
func (h *PaymentHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID   string  `json:"user_id"`
		Amount   float64 `json:"amount"`
		Currency string  `json:"currency"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	payment, err := h.svc.CreatePayment(r.Context(), req.UserID, req.Amount, req.Currency)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(payment)
}

// GetPayment godoc
// @Summary Get payment by ID
// @Description Retrieve a payment by its ID
// @Tags payments
// @Accept json
// @Produce json
// @Param id query string true "Payment ID"
// @Success 200 {object} models.PaymentResponse "Payment found"
// @Failure 400 {object} models.ErrorResponse "Missing ID"
// @Failure 404 {object} models.ErrorResponse "Payment not found"
// @Router /payments [get]
func (h *PaymentHandler) GetPayment(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id required", http.StatusBadRequest)
		return
	}

	payment, err := h.svc.GetPaymentByID(r.Context(), id)
	if err != nil {
		http.Error(w, "payment not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(payment)
}

// GetPaymentsByUserID godoc
// @Summary Get payments by user ID
// @Description Retrieve payments for a specific user, returns last payment status
// @Tags payments
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} models.GetPaymentsByUserResponse "User payments retrieved"
// @Failure 404 {object} models.ErrorResponse "Payments not found"
// @Router /payments/user/{user_id} [get]
func (h *PaymentHandler) GetPaymentsByUserID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	payments, err := h.svc.GetPaymentsByUser(r.Context(), userID)
	if err != nil {
		http.Error(w, "payments not found", http.StatusNotFound)
		return
	}

	var status string
	if len(payments) > 0 {
		status = string(payments[0].Status)
	} else {
		status = "none"
	}

	json.NewEncoder(w).Encode(struct {
		Status string `json:"status"`
	}{
		Status: status,
	})
}

// GetPaymentsByUserID godoc
// @Summary Get payments by user ID
// @Description Retrieve payments for a specific user, returns last payment status
// @Tags payments
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} models.GetPaymentsByUserResponse "User payments retrieved"
// @Failure 404 {object} models.ErrorResponse "Payments not found"
// @Router /payments/user/{user_id} [get]
func (h *PaymentHandler) UpdatePaymentStatus(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PaymentID string               `json:"payment_id"`
		Status    models.PaymentStatus `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	if err := h.svc.UpdatePaymentStatus(r.Context(), req.PaymentID, req.Status); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
