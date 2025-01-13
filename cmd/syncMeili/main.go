package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"

	"github.com/Anuolu-2020/sync-meili/internal/config"
	"github.com/Anuolu-2020/sync-meili/internal/db"
	"github.com/Anuolu-2020/sync-meili/internal/sync"
	"github.com/Anuolu-2020/sync-meili/internal/webhook"
	"github.com/Anuolu-2020/sync-meili/pkg/env"
)

func main() {
	config, err := config.LoadConfig("./config.yaml")
	if err != nil {
		fmt.Printf("Failed to read sync-meili config: %v", err)
		os.Exit(1)
	}

	// Initialize database
	db.InitDB(*config)

	batchChannel := make(chan sync.MeiliSearchRequest)

	if config.Database.Type == "postgres" {
		// Initialize Syncing
		sync.CreateTrigger(*config)

		// Start syncing
		go sync.ListenAndSync(config, batchChannel)
		log.Print("Started listening for database updates")

		defer db.DBManager.Conn.Close()

	} else {
		go sync.RegisterAndStartCanal(config, batchChannel)
		log.Print("Started listening for database updates")

		defer db.DBManager.MysqlCanal.Close()
	}

	// Sync documents in batches
	go sync.SyncInBatchRequest(batchChannel, config)

	// Webhook endpoint
	http.HandleFunc("/webhook", webhook.WebhookHandler)

	// start server for webhook
	log.Printf("Webhook server listening on http://localhost%s", env.GetEnv("PORT"))
	http.ListenAndServe(env.GetEnv("PORT"), nil)
}
