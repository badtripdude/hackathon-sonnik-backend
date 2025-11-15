package auth

import (
	"crypto/rsa"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GenerateAccessToken(userID, email string, privKey *rsa.PrivateKey, ttl time.Duration) (string, time.Time, error) {
	if privKey == nil {
		return "", time.Time{}, errors.New("private key is required")
	}

	exp := time.Now().UTC().Add(ttl)

	claims := &jwt.MapClaims{
		"sub":   userID,
		"exp":   jwt.NewNumericDate(exp),
		"iat":   jwt.NewNumericDate(time.Now().UTC()),
		"jti":   uuid.New().String(),
		"email": email,
	}

	rawToken := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	signedToken, err := rawToken.SignedString(privKey)
	if err != nil {
		return "", time.Time{}, err
	}

	return signedToken, exp, nil
}
