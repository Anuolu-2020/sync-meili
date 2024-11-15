package db

import (
	"database/sql"
	"log"

	"github.com/go-mysql-org/go-mysql/canal"

	"github.com/Anuolu-2020/sync-meili/internal/config"
	"github.com/Anuolu-2020/sync-meili/internal/db/mysql"
	"github.com/Anuolu-2020/sync-meili/internal/db/postgres"
	"github.com/Anuolu-2020/sync-meili/pkg/env"
)

type DBConnection struct {
	Conn       *sql.DB
	MysqlCanal *canal.Canal
}

var DBManager = &DBConnection{}

func InitDB(config config.SyncMeiliConfig) {
	switch config.Database.Type {
	case "postgres":
		conn, err := postgres.Connect(env.GetEnv("DB_CONNECTION_STRING"))
		if err != nil {
			log.Fatalf("Failed to connect to Postgres: %v", err)
		}

		DBManager.Conn = conn.Conn
	case "mysql":
		mysqlCanal, err := mysql.Connect(config)
		if err != nil {
			log.Fatalf("Failed to connect to Mysql: %v", err)
		}

		DBManager.MysqlCanal = mysqlCanal.Canal

	default:
		log.Fatalf("Database type not supported")
	}
}
