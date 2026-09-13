package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/xenptr/go-projects/ecommerce-api/internal/config"
	"github.com/xenptr/go-projects/ecommerce-api/internal/database"
	"github.com/xenptr/go-projects/ecommerce-api/internal/handler"
	"github.com/xenptr/go-projects/ecommerce-api/internal/repository"
	"github.com/xenptr/go-projects/ecommerce-api/internal/routes"
)

var shutdownTimeout = 10 * time.Second

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config.Load()

	pool, err := database.Open(cfg)
	if err != nil {
		log.Fatal("error establishing connection")
	}
	defer pool.Close()

	repo := repository.New(pool)

	mux := http.NewServeMux()
	h := handler.New(repo)

	routes.RegisterRoutes(mux, h)

	srv := &http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: mux,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)

	serverErr := make(chan error, 1)
	go func() {
		log.Printf("server started on %s", cfg.AppPort)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	select {
	case <-quit:
		log.Println("performing graceful shutdown")
	case err := <-serverErr:
		log.Printf("server terminated unexpectedly: %v", err)
		return
	}

	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("server shutdown error: %v", err)
		return
	}

	log.Println("server stopped!")
}
