package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"github.com/badtripdude/hackathon-sonnik-backend/services/user/internal/models"
	"time"

	"github.com/google/uuid"
)

func GenerateRefreshToken(userID string, ttl time.Duration) (*models.RefreshToken, string, error) {
	plain, err := generatePlainToken(64)
	if err != nil {
		return nil, "", err
	}

	hashed := HashToken(plain)

	token := &models.RefreshToken{
		ID:        uuid.New().String(),
		UserID:    userID,
		TokenHash: hashed,
		ExpiresAt: time.Now().UTC().Add(ttl),
		CreatedAt: time.Now().UTC(),
	}

	return token, plain, nil
}

func generatePlainToken(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func HashToken(plain string) string {
	h := sha256.Sum256([]byte(plain))
	return hex.EncodeToString(h[:])
}
