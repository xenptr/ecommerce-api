package handler

import (
	"github.com/xenptr/ecommerce-api/internal/repository"
	"github.com/xenptr/ecommerce-api/internal/service"
)

type Handler struct {
	userRepo    repository.UserRepository
	authService *service.AuthService
}

func New(store repository.Store, authService *service.AuthService) *Handler {
	return &Handler{
		userRepo:    store,
		authService: authService,
	}
}
