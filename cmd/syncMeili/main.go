package main

import (
	"log"

	"github.com/Anuolu-2020/sync-meili/internal/config"
	"github.com/Anuolu-2020/sync-meili/internal/db"
)

func main() {
	config, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to read sync-meili config: %v", err)
	}

	// Initialize database
	db.InitDB(config)

	defer db.DBManager.Conn.Close()
}
