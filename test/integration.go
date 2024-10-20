package test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/Anuolu-2020/sync-meili/internal/config"
	"github.com/Anuolu-2020/sync-meili/internal/db"
)

func TestPostgresAndMeilisearch(t *testing.T) {
	ctx := context.Background()

	postgresContainer, err := postgres.Run(
		ctx,
		"postgres:16-alpine",
		postgres.WithUsername("user"),
		postgres.WithDatabase("test"),
		postgres.WithPassword("password"),
	)
	if err != nil {
		t.Fatalf("Postgres Container failed to start: %v", err)
	}

	postgresConnectionString, err := postgresContainer.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("Failed to get postgres container connection string: %v", err)
	}

	os.Setenv("DB_URL", postgresConnectionString)

	config, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to read sync-meili config: %v", err)
	}

	// Initialize database
	db.InitDB(config)
}
