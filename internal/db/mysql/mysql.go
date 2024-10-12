package mysql

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type MySqlConnection struct {
	Conn *sql.DB
}

func Connect(url string) (*MySqlConnection, error) {
	conn, err := sql.Open("mysql", url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return nil, err
	}

	return &MySqlConnection{conn}, nil
}

func (db *MySqlConnection) Close() {
	db.Conn.Close()
}
