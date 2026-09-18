package routes

import (
	"net/http"

	"github.com/xenptr/ecommerce-api/internal/handler"
)

func RegisterRoutes(m *http.ServeMux, h *handler.Handler) {
	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Welcome to the E-Commerce API!"))
	})

	m.HandleFunc("POST /api/v1/auth/register", h.Register)
	m.HandleFunc("POST /api/v1/auth/login", h.Login)
}
