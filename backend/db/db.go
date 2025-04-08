package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"grpcShop.com/backend/config"
)

func NewDB() *sql.DB {
	// Load config
	config := config.LoadConfig()
	fmt.Println("Config loaded:", config)
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
