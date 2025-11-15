package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/badtripdude/hackathon-sonnik-backend/services/user/internal/service"
	"github.com/gorilla/mux"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// Register godoc
// @Summary Register new user
// @Description Register a user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Register request body"
// @Success 200 {object} models.TokenPair "User tokens"
// @Failure 400 {object} models.ErrorResponse "Invalid body or user already exists"
// @Router /auth/register [post]
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email     string  `json:"email"`
		Password  string  `json:"password"`
		Username  *string `json:"username"`
		BirthDate *string `json:"birth_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var birthDatePtr *time.Time
	if req.BirthDate != nil {
		parsed, err := time.Parse("2006-01-02", *req.BirthDate)
		if err != nil {
			http.Error(w, "invalid birth_date format", http.StatusBadRequest)
			return
		}
		birthDatePtr = &parsed
	}

	tokens, err := h.svc.Register(r.Context(), req.Email, req.Password, req.Username, birthDatePtr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(tokens)
}

// UpdateUser godoc
// @Summary Update user data
// @Description Update email or username for user
// @Tags user
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body models.UpdateUserRequest true "Fields to update"
// @Success 200 {object} models.User "Updated user"
// @Failure 400 {object} models.ErrorResponse "Invalid request"
// @Failure 404 {object} models.ErrorResponse "User not found"
// @Router /auth/user/{id} [patch]
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    *string `json:"email"`
		Username *string `json:"username"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	user, err := h.svc.UpdateUser(r.Context(), id, body.Email, body.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login request body"
// @Success 200 {object} models.TokenPair "Successful login"
// @Failure 400 {object} models.ErrorResponse "Invalid body, wrong password, or user not found"
// @Router /auth/login [post]
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tokens, err := h.svc.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(tokens)
}

// Refresh godoc
// @Summary Refresh tokens
// @Description Generate new access and refresh tokens using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RefreshRequest true "Refresh request body"
// @Success 200 {object} models.TokenPair "Refreshed tokens"
// @Failure 400 {object} models.ErrorResponse "Invalid request body"
// @Failure 401 {object} models.ErrorResponse "Refresh token invalid, expired, or revoked"
// @Router /auth/refresh [post]
func (h *UserHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tokens, err := h.svc.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	json.NewEncoder(w).Encode(tokens)
}

// Logout godoc
// @Summary Logout user
// @Description Invalidate refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RefreshRequest true "Logout request body"
// @Success 204 "No content"
// @Failure 400 {object} models.ErrorResponse "Invalid body or token not found"
// @Router /auth/logout [post]
func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.svc.Logout(r.Context(), req.RefreshToken); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Me godoc
// @Summary Get current user info
// @Description Return data of the authenticated user
// @Tags user
// @Produce json
// @Success 200 {object} models.User "User info"
// @Failure 404 {object} models.ErrorResponse "User not found"
// @Router /auth/me [get]
func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	user, err := h.svc.GetByID(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(user)
}
