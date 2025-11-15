package service

import (
	"context"
	"crypto/rsa"
	"errors"
	"time"

	"github.com/badtripdude/hackathon-sonnik-backend/services/user/internal/auth"
	"github.com/badtripdude/hackathon-sonnik-backend/services/user/internal/config"
	"github.com/badtripdude/hackathon-sonnik-backend/services/user/internal/models"
	errs "github.com/badtripdude/hackathon-sonnik-backend/services/user/pkg/errors"
)

type RefreshTokenRepo interface {
	Save(ctx context.Context, token *models.RefreshToken) error
	GetByTokenHash(ctx context.Context, token string) (*models.RefreshToken, error)
	RevokeByTokenHash(ctx context.Context, token string) error
	RevokeAllByUserId(ctx context.Context, userId string) error
}

type UserRepo interface {
	Save(ctx context.Context, user *models.User) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByID(ctx context.Context, id string) (*models.User, error)
	DeleteByID(ctx context.Context, id string) error
	Update(ctx context.Context, id string, email *string, username *string) (*models.User, error)
	UpdateAvatar(ctx context.Context, id string, avatar []byte) error
}

type UserService struct {
	cfg              config.Config
	userRepo         UserRepo
	refreshTokenRepo RefreshTokenRepo
	privateKey       *rsa.PrivateKey
}

func NewUserService(cfg config.Config, userRepo UserRepo, refreshTokenRepo RefreshTokenRepo, key *rsa.PrivateKey) *UserService {
	return &UserService{
		cfg:              cfg,
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		privateKey:       key,
	}
}

func (s *UserService) Register(ctx context.Context, email, password string, username *string, birthDate *time.Time) (*models.TokenPair, error) {
	_, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil {
		return nil, errs.ErrUserAlreadyExists
	}
	if !errors.Is(err, errs.ErrNotFound) {
		return nil, err
	}

	hashedPassword, err := auth.HashPassword(password, s.cfg.BcryptCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:        email,
		PasswordHash: hashedPassword,
		Username:     username,
		BirthDate:    birthDate,
		CreatedAt:    time.Now(),
	}

	if err := s.userRepo.Save(ctx, user); err != nil {
		return nil, err
	}

	refreshToken, plain, err := auth.GenerateRefreshToken(user.ID, s.cfg.RefreshTokenTTL)
	if err != nil {
		return nil, err
	}

	if err := s.refreshTokenRepo.Save(ctx, refreshToken); err != nil {
		return nil, err
	}

	accessToken, exp, err := auth.GenerateAccessToken(user.ID, email, s.privateKey, s.cfg.AccessTokenTTL)
	if err != nil {
		return nil, err
	}

	return &models.TokenPair{
		UserID:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: plain,
		ExpiresAt:    exp,
	}, nil
}

func (s *UserService) Login(ctx context.Context, email, password string) (*models.LoginResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}

	if !auth.CheckPasswordHash(password, user.PasswordHash) {
		return nil, errs.ErrPasswordMismatch
	}

	refreshToken, plain, err := auth.GenerateRefreshToken(user.ID, s.cfg.RefreshTokenTTL)
	if err != nil {
		return nil, err
	}

	_ = s.refreshTokenRepo.RevokeAllByUserId(ctx, user.ID)

	if err := s.refreshTokenRepo.Save(ctx, refreshToken); err != nil {
		return nil, err
	}

	accessToken, exp, err := auth.GenerateAccessToken(user.ID, email, s.privateKey, s.cfg.AccessTokenTTL)
	if err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		UserID:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: plain,
		ExpiresAt:    exp,
	}, nil
}

func (s *UserService) Logout(ctx context.Context, plainRefreshToken string) error {
	hashed := auth.HashToken(plainRefreshToken)
	err := s.refreshTokenRepo.RevokeByTokenHash(ctx, hashed)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return errs.ErrRefreshTokenNotFound
		}
		return err
	}
	return nil
}

func (s *UserService) GetByID(ctx context.Context, id string) (*models.User, error) {
	var user *models.User

	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id string, email *string, username *string) (*models.User, error) {
	user, err := s.userRepo.Update(ctx, id, email, username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Refresh(ctx context.Context, plainRefreshToken string) (*models.TokenPair, error) {
	hashed := auth.HashToken(plainRefreshToken)

	rt, err := s.refreshTokenRepo.GetByTokenHash(ctx, hashed)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return nil, errs.ErrRefreshTokenNotFound
		}
		return nil, err
	}

	if rt.RevokedAt != nil {
		return nil, errs.ErrRefreshTokenRevoked
	}
	if rt.ExpiresAt.Before(time.Now()) {
		return nil, errs.ErrRefreshTokenExpired
	}

	user, err := s.userRepo.GetByID(ctx, rt.UserID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}

	newRT, newPlain, err := auth.GenerateRefreshToken(user.ID, s.cfg.RefreshTokenTTL)
	if err != nil {
		return nil, err
	}

	if err := s.refreshTokenRepo.Save(ctx, newRT); err != nil {
		return nil, err
	}

	if err := s.refreshTokenRepo.RevokeByTokenHash(ctx, hashed); err != nil {
		return nil, err
	}

	accessToken, exp, err := auth.GenerateAccessToken(user.ID, user.Email, s.privateKey, s.cfg.AccessTokenTTL)
	if err != nil {
		return nil, err
	}

	return &models.TokenPair{
		UserID:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: newPlain,
		ExpiresAt:    exp,
	}, nil
}

func (s *UserService) UploadAvatar(ctx context.Context, id string, avatar []byte) error {
	return s.userRepo.UpdateAvatar(ctx, id, avatar)
}

func (s *UserService) UpdateAvatar(ctx context.Context, id string, avatar []byte) error {
	return s.userRepo.UpdateAvatar(ctx, id, avatar)
}
