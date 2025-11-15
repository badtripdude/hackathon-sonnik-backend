package server

import (
	"net/http"

	"github.com/badtripdude/hackathon-sonnik-backend/services/robokassa/internal/controller/http/handlers"
	"github.com/badtripdude/hackathon-sonnik-backend/services/robokassa/internal/service"
	"github.com/gorilla/mux"
)

func NewRouter(paymentSvc *service.PaymentService) http.Handler {
	r := mux.NewRouter()
	h := handlers.NewPaymentHandler(paymentSvc)

	r.HandleFunc("/payments", h.CreatePayment).Methods("POST")
	r.HandleFunc("/payments", h.GetPayment).Methods("GET")
	r.HandleFunc("/payments/status", h.UpdatePaymentStatus).Methods("PATCH")
	r.HandleFunc("/payments/user/{user_id}", h.GetPaymentsByUserID).Methods("GET")

	return r
}
