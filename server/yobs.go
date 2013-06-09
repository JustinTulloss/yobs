package main

import (
	"fmt"
    _ "github.com/lib/pq"
    "database/sql"
    "net/http"
    "os"
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
}

func NewUser(email string) *User {
	fmt.Printf("Creating user %s\n", email)
	u := new(User)
	u.email = email
	return u
}

func InsertUser(user *User, db *sql.DB) *User {
	fmt.Printf("Inserting user with email %s\n", user.email)
	stmt, err := db.Prepare("INSERT INTO users (email) VALUES ($1) RETURNING id;")
	var id int64
	err = stmt.QueryRow(user.email).Scan(&id)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	user.id = id
	fmt.Printf("New user has id %d\n", id)
	return user
}

func UserCount(db *sql.DB) int {
	var count int
	db.QueryRow("SELECT COUNT(*) FROM users;").Scan(&count)
	return count
}

func Users(db *sql.DB) ([]*User, error) {
	fmt.Printf("Querying for users.\n")
	rows, err := db.Query("SELECT id, email FROM users;")
	if err != nil {
		fmt.Printf("Error querying for users:\n%s\n", err)
		return nil, err
	}
	count := UserCount(db)
	fmt.Printf("Found %d users.\n", count)
	var users []*User
	for rows.Next() {
		var id int64
		var email string
		err = rows.Scan(&id, &email)
		user := new(User)
		user.id = id
		user.email = email
		users = append(users, user)
	}
	return users, nil
}

func main() {
	http.HandleFunc("/users", users)
	http.HandleFunc("/users/new", new_user)

	fmt.Printf("Listening...")
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
    if err != nil {
		panic(err)
    }
}

func new_user(res http.ResponseWriter, req *http.Request) {
	db, _ := initDB()
	params := req.URL.Query()
	email := params["email"][0]
	
	user := new(User)
	user.email = email
	user = InsertUser(user, db)

	fmt.Fprintf(res, "%s", user)
	defer db.Close()
}

func users(res http.ResponseWriter, req *http.Request) {
	db, err := initDB()
	if err != nil {
		fmt.Fprintf(res, "There was an error querying: %s\n", err)
	}
	users, _ := Users(db)
	for i :=0; i < len(users); i++ {
		fmt.Fprintf(res, "%d: %s\n", users[i].id, users[i].email)
	}
	
	defer db.Close()
}