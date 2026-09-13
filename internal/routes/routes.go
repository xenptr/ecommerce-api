package routes

import (
	"net/http"

	"github.com/xenptr/go-projects/ecommerce-api/internal/handler"
)

func RegisterRoutes(m *http.ServeMux, h *handler.Handler) {
	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Welcome to the E-Commerce API!"))
	})
}
