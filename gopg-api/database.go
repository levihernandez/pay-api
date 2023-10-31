package main

import (
	"fmt"
	"github.com/go-pg/pg/v10"
	"log"
)

type DBConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database   string `json:"database"`
}

func InitDB(config *DBConfig) *pg.DB {
	db := pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		User:     config.User,
		Password: config.Password,
		Database: config.Database,
	})

	// Check if the database is reachable
	_, err := db.Exec("SELECT 1")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	return db
}
