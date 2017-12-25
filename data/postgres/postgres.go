package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/lib/pq" // blank import for PG
)

// Postgres is a DB wrapper
type Postgres struct {
	*sql.DB
}

// NewPostgres factory
func NewPostgres(user, password, dbName, host, port, maxIdle, maxOpen string) *Postgres {
	// user = os.Getenv("PG_DB_USER")
	// password = os.Getenv("PG_PASSWORD")
	// dbName = os.Getenv("PG_DB_NAME")
	// host = os.Getenv("PG_HOST")
	// port = os.Getenv("PG_PORT")

	if user == "" || password == "" || dbName == "" || host == "" || port == "" {
		log.Printf("Invalid vars in NewPostgres: \n PG_DB_USER:%s \nPG_PASSWORD:%s \nPG_DB_NAME: %s \nPG_HOST: %s \nPG_PORT: %s\n",
			user, password, dbName, host, port)
		return nil
	}

	db, err := sql.Open("postgres",
		fmt.Sprintf("user=%s password='%s' dbname=%s host=%s port=%s sslmode=disable",
			user,
			password,
			dbName,
			host,
			port,
		),
	)
	if err != nil {
		log.Fatalf("[NewPG] %s\n", err.Error())
	}
	maxIdle = os.Getenv("PG_MAX_IDLE")
	if maxIdle == "" {
		maxIdle = "10"
	}
	maxIdleConnections, err := strconv.Atoi(maxIdle)
	if err != nil {
		maxIdleConnections = 10
	}
	maxOpen = os.Getenv("PG_MAX_OPEN")
	if maxOpen == "" {
		maxOpen = "10"
	}

	maxOpenConnections, err := strconv.Atoi(maxOpen)
	if err != nil {
		maxOpenConnections = 10
	}

	db.SetMaxIdleConns(maxIdleConnections)
	db.SetMaxOpenConns(maxOpenConnections)
	return &Postgres{db}
}
