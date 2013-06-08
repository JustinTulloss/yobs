package main

import (
	"fmt"
    _ "github.com/lib/pq"
    "database/sql"
	"time"
)

func initDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", "dbname=yobs sslmode=disable")
	if err != nil {
		fmt.Printf("Error opening DB connection:\n%s\n", err)
		return db, err
	} else {
		fmt.Printf("Connected to database.\n")
	}
	return db, err
}

type User struct {
	email string
	id int64
	created int64 // seconds
}

func NewUser(email string) *User {
	fmt.Printf("Creating user %s\n", email)
	u := new(User)
	u.email = email
	u.created = time.Now().UTC().Unix()
	return u
}

func Users(db *sql.DB) ([]User, error) {
	fmt.Printf("Querying for users.\n")
	rows, err := db.Query("SELECT id, email FROM users;")
	if err != nil {
		fmt.Printf("Error querying for users:\n%s\n", err)
		return nil, err
	}
	
	for rows.Next() {
		var id int
		var email string
		err = rows.Scan(&id, &email)
		fmt.Printf("%d: %s\n", id, email)
	}
	return nil, nil
}

func main() {
	db, err := initDB()
	if err != nil {
		fmt.Printf("There was an error querying: %s\n", err)
	}
	Users(db)
	defer db.Close()
}
