package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type PostgresConnection struct {
	Conn *sql.DB
}

func Connect(url string) (*PostgresConnection, error) {
	conn, err := sql.Open("postgres", url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return nil, err
	}

	return &PostgresConnection{conn}, nil
}

func (db *PostgresConnection) Close() {
	log.Print("closing postgres connection")
	db.Conn.Close()
}
