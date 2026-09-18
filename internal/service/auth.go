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
}

func NewAuthService(userRepo repository.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
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

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) error {
	validate := validator.New()

	if err := validate.Struct(req); err != nil {
		return err
	}
	user, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return err
	}

	if err := auth.CheckPassword(req.Password, user.PasswordHash); err != nil {
		return err
	}

	return nil
}
