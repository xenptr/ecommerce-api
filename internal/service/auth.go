package service

import (
	"context"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/xenptr/ecommerce-api/internal/auth"
	"github.com/xenptr/ecommerce-api/internal/dto"
	"github.com/xenptr/ecommerce-api/internal/models"
	"github.com/xenptr/ecommerce-api/internal/repository"
	"github.com/xenptr/ecommerce-api/internal/session"
	"github.com/xenptr/ecommerce-api/internal/token"
)

type AuthService struct {
	userRepo     repository.UserRepository
	tokens       AuthTokens
	refreshStore session.RefreshTokenStore
}

type AuthTokens interface {
	Generate(userID int64) (string, error)
	GenerateRefresh(userID int64) (string, error)
	ParseRefresh(tokenString string) (token.RefreshTokenInfo, error)
}

func NewAuthService(userRepo repository.Store, tokens AuthTokens, refreshStore session.RefreshTokenStore) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		tokens:       tokens,
		refreshStore: refreshStore,
	}
}

func (s *AuthService) Register(ctx context.Context, req dto.RegisterRequest) error {
	validate := validator.New()

	if err := validate.Struct(req); err != nil {
		return err
	}

	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		return err
	}

	user := models.User{
		Email:        req.Email,
		PasswordHash: passwordHash,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Role:         "customer",
		IsActive:     true,
	}

	return s.userRepo.CreateUser(ctx, user)
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (dto.AuthResponse, error) {
	var resp dto.AuthResponse

	validate := validator.New()

	if err := validate.Struct(req); err != nil {
		return resp, err
	}

	user, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return resp, err
	}

	if err := auth.CheckPassword(req.Password, user.PasswordHash); err != nil {
		return resp, err
	}

	return s.issueTokenPair(ctx, user.ID)
}

func (s *AuthService) RefreshToken(ctx context.Context, req dto.RefreshRequest) (dto.AuthResponse, error) {
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return dto.AuthResponse{}, err
	}

	tokenInfo, err := s.tokens.ParseRefresh(req.RefreshToken)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	exists, err := s.refreshStore.Exists(ctx, tokenInfo.UserID, tokenInfo.JTI)
	if err != nil {
		return dto.AuthResponse{}, err
	}
	if !exists {
		return dto.AuthResponse{}, errors.New("refresh token is invalid or expired")
	}

	// Verify the user account still exists.
	if _, err = s.userRepo.GetUserByID(ctx, tokenInfo.UserID); err != nil {
		return dto.AuthResponse{}, err
	}

	if err := s.refreshStore.Revoke(ctx, tokenInfo.UserID, tokenInfo.JTI); err != nil {
		return dto.AuthResponse{}, err
	}

	return s.issueTokenPair(ctx, tokenInfo.UserID)
}

func (s *AuthService) issueTokenPair(ctx context.Context, userID int64) (dto.AuthResponse, error) {
	accessToken, err := s.tokens.Generate(userID)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	refreshToken, err := s.tokens.GenerateRefresh(userID)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	tokenInfo, err := s.tokens.ParseRefresh(refreshToken)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	if err := s.refreshStore.Save(ctx, userID, tokenInfo.JTI, token.RefreshTokenTTL); err != nil {
		return dto.AuthResponse{}, err
	}

	return dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
