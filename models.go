package main

import (
    _ "github.com/lib/pq"
    "database/sql"
	"log"
)

// Database Interaction
func initDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", "dbname=yobs sslmode=disable")
	if err != nil {
		log.Printf("Error opening DB connection:\n%s\n", err)
		return db, err
	} 
	return db, err
}
