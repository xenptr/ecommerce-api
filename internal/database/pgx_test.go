package database

import (
	"context"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/xenptr/go-projects/ecommerce-api/internal/config"
)

func startPostgresContainer() (
	context.Context,
	testcontainers.Container,
	*config.Config,
	error,
) {
	ctx := context.Background()

	container, err := postgres.Run(
		ctx,
		"postgres:latest",
		postgres.WithDatabase("database"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		return nil, nil, nil, err
	}

	host, err := container.Host(ctx)
	if err != nil {
		container.Terminate(ctx)
		return nil, nil, nil, err
	}

	port, err := container.MappedPort(ctx, "5432/tcp")
	if err != nil {
		container.Terminate(ctx)
		return nil, nil, nil, err
	}

	cfg := &config.Config{
		DBUser: "user",
		DBPass: "password",
		DBHost: host,
		DBPort: port.Port(),
		DBName: "database",
	}

	return ctx, container, cfg, nil
}

func TestOpen(t *testing.T) {
	ctx, container, cfg, err := startPostgresContainer()
	if err != nil {
		t.Fatalf("failed to start postgres: %v", err)
	}
	defer container.Terminate(ctx)

	pool, err := Open(cfg)
	if err != nil {
		t.Fatalf("Open() failed: %v", err)
	}
	defer pool.Close()

	if pool == nil {
		t.Fatal("expected pool, got nil")
	}
}
