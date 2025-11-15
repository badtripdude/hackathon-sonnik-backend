package auth

import (
	"crypto/rsa"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func LoadPublicKey(path string) (*rsa.PublicKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	pub, err := jwt.ParseRSAPublicKeyFromPEM(data)
	if err != nil {
		return nil, err
	}

	return pub, nil
}
