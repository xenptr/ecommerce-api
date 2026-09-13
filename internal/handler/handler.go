package handler

import "github.com/xenptr/go-projects/ecommerce-api/internal/repository"

type Handler struct {
	userRepo repository.UserRepository
}

func New(store repository.Store) *Handler {
	return &Handler{
		userRepo: store,
	}
}
