package server

import (
	"crypto/rsa"
	"github.com/badtripdude/hackathon-sonnik-backend/services/user/internal/controller/http/handlers"
	"github.com/badtripdude/hackathon-sonnik-backend/services/user/internal/service"
	"github.com/badtripdude/hackathon-sonnik-backend/services/user/middleware"
	"github.com/gorilla/mux"
	"net/http"
)

func NewRouter(userSvc *service.UserService, pubKey *rsa.PublicKey) http.Handler {
	r := mux.NewRouter()
	h := handlers.NewUserHandler(userSvc)

	r.HandleFunc("/auth/register", h.Register).Methods("POST")
	r.HandleFunc("/auth/login", h.Login).Methods("POST")
	r.HandleFunc("/auth/refresh", h.Refresh).Methods("POST")
	r.HandleFunc("/auth/logout", h.Logout).Methods("POST")
	r.Handle("/auth/me", middleware.JWTMiddleware(pubKey)(h.Me)).Methods("GET")

	return r
}
