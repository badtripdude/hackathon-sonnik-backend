package errors

import "errors"

var (
	// repo errs
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")

	// access token errs
	ErrAccessTokenExpired = errors.New("access token has been expires")
	ErrAccessTokenInvalid = errors.New("access token is invalid")

	// refresh token errs
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
	ErrRefreshTokenExpired  = errors.New("refresh token has expired")
	ErrRefreshTokenRevoked  = errors.New("refresh token has been revoked")

	// user errs
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")

	// pswd err
	ErrPasswordMismatch = errors.New("password mismatch")
)
