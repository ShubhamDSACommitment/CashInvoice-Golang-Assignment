package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/CashInvoice-Golang-Assignment/internal/config"
	_ "github.com/go-sql-driver/mysql"
)

func Connect(cfg *config.Config) *sql.DB {

	host := cfg.DBHost
	port := cfg.DBPort
	user := cfg.DBUser
	pass := cfg.DBPass
	name := cfg.DBName
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pass, host, port, name)

	log.Println("Connecting to database...", dsn)
	var db *sql.DB
	var err error

	db, err = sql.Open("mysql", dsn)
	if err == nil && db.Ping() == nil {
		return db
	}
	log.Println(err)
	panic("Could not connect to MySQL")
}

func RunMigrations(db *sql.DB) {
	queries := []string{
		`
        CREATE TABLE IF NOT EXISTS users (
            id VARCHAR(36) PRIMARY KEY,
            email VARCHAR(255) UNIQUE,
            password TEXT,
            role VARCHAR(20)
        );
        `,
		`
        CREATE TABLE IF NOT EXISTS tasks (
            id VARCHAR(36) PRIMARY KEY,
            title VARCHAR(255) NOT NULL,
            description TEXT,
            status VARCHAR(20) NOT NULL,
            user_id VARCHAR(36) NOT NULL,
            created_at TIMESTAMP NOT NULL,
            updated_at TIMESTAMP NOT NULL
        );
        `,
	}

	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			log.Println("Migration error:", err)
		}
	}
}
