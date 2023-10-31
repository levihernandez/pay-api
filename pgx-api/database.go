package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

func initDB(config *DBConfig) (*pgxpool.Pool, error) {
	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
		config.User, config.Password, config.Host, config.Port, config.Database)

	dbpool, err := pgxpool.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
		return nil, err
	}

	err = dbpool.Ping(context.Background())
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
		return nil, err
	}

	log.Println("Connection to database successful.")
	return dbpool, nil
}
