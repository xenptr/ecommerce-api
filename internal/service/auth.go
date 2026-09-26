package service

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/xenptr/ecommerce-api/internal/auth"
	"github.com/xenptr/ecommerce-api/internal/dto"
	"github.com/xenptr/ecommerce-api/internal/models"
	"github.com/xenptr/ecommerce-api/internal/repository"
)

type AuthService struct {
	userRepo repository.UserRepository
	tokens   AuthTokens
}

type AuthTokens interface {
	Generate(userID int64) (string, error)
	GenerateRefresh(userID int64) (string, error)
	ParseRefresh(tokenString string) (int64, error)
}

func NewAuthService(userRepo repository.Store, tokens AuthTokens) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		tokens:   tokens,
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

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (string, error) {
	validate := validator.New()

	if err := validate.Struct(req); err != nil {
		return "", err
	}
	user, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return "", err
	}

	if err := auth.CheckPassword(req.Password, user.PasswordHash); err != nil {
		return "", err
	}

	token, err := s.tokens.Generate(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, req dto.RefreshRequest) (dto.AuthResponse, error) {
	var resp dto.AuthResponse

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return resp, err
	}

}
