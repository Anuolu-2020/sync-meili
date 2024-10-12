package db

import (
	"database/sql"
	"log"

	"github.com/Anuolu-2020/sync-meili/internal/config"
	"github.com/Anuolu-2020/sync-meili/internal/db/mysql"
	"github.com/Anuolu-2020/sync-meili/internal/db/postgres"
)

type DBConnection struct {
	Conn *sql.DB
}

var DBManager *DBConnection

func InitDB(config *config.SyncMeiliConfig) {
	switch config.Database.Type {
	case "postgres":
		conn, err := postgres.Connect("hello")
		if err != nil {
			log.Fatalf("Failed to connect to Postgres: %v", err)
		}

		DBManager.Conn = conn.Conn
	case "mysql":
		conn, err := mysql.Connect("hello")
		if err != nil {
			log.Fatalf("Failed to connect to Mysql: %v", err)
		}

		DBManager.Conn = conn.Conn
	default:
		log.Fatalf("Database type not supported")
	}
}
