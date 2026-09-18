package routes

import (
	"net/http"

	"github.com/xenptr/ecommerce-api/internal/handler"
	"github.com/xenptr/ecommerce-api/internal/middleware"
	"github.com/xenptr/ecommerce-api/internal/token"
)

func RegisterRoutes(m *http.ServeMux, h *handler.Handler, parser token.Parser) {
	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Welcome to the E-Commerce API!"))
	})

	m.HandleFunc("POST /api/v1/auth/register", h.Register)
	m.HandleFunc("POST /api/v1/auth/login", h.Login)

	m.Handle("GET /api/v1/auth",
		middleware.Auth(parser)(
			http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					_, _ = w.Write([]byte("You are authenticated!"))
				}),
		),
	)
}
