package models

type RegisterRequest struct {
    Email     string `json:"email" example:"test@example.com"`
    Password  string `json:"password" example:"password123"`
    Username  string `json:"username" example:"JohnDoe"`
    BirthDate string `json:"birth_date" example:"1990-01-01"`
}

type LoginRequest struct {
    Email    string `json:"email" example:"test@example.com"`
    Password string `json:"password" example:"password123"`
}

type RefreshRequest struct {
    RefreshToken string `json:"refresh_token"`
}

type ErrorResponse struct {
    Error string `json:"error" example:"password mismatch"`
}
