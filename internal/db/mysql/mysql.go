package mysql

import (
	"fmt"
	"os"

	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/siddontang/go-log/log"

	"github.com/Anuolu-2020/sync-meili/internal/config"
	"github.com/Anuolu-2020/sync-meili/pkg/env"
)

type MySqlConnection struct {
	Canal *canal.Canal
}

func Connect(config config.SyncMeiliConfig) (*MySqlConnection, error) {
	cfg := canal.NewDefaultConfig()
	cfg.Addr = env.GetEnv("MYSQL_ADDR")
	cfg.User = env.GetEnv("MYSQL_USER")
	cfg.Password = env.GetEnv("MYSQL_PASSWORD")
	cfg.Dump.TableDB = env.GetEnv("MYSQL_TABLE_DB")
	// databaseTables := make([]string, 0)

	// disable embedded canal logs
	// cfg.Logger = log.NewDefault(&log.NullHandler{})
	// var db string
	// dbs := make([]string, 0)

	// for _, mapping := range config.Sync.Mappings {
	// 	// db = strings.Split(mapping.DatabaseTable, ".")[0]
	// 	// dbs = append(dbs, db)
	// 	databaseTables = append(databaseTables, mapping.DatabaseTable)
	// }
	//
	// cfg.Dump.Tables = databaseTables

	c, err := canal.NewCanal(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return nil, err
	}

	return &MySqlConnection{Canal: c}, nil
}

func (canal MySqlConnection) Close() {
	log.Print("closing mysql canal")
	canal.Canal.Close()
}
