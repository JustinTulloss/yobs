package main

import (
	"fmt"
    _ "github.com/lib/pq"
    "database/sql"
    "net/http"
    "os"
	"strconv"
	"encoding/json"
)

// Models
type User struct {
	Facebook_id int64
	Id int64
}

func (u User) Transactions() *TransactionCollection {
	db, _ := initDB()
	stmt, _ := db.Prepare("SELECT id, owner_id, amount, description FROM transactions WHERE owner_id = '$1';")
	rows, _ := stmt.Query(u.Id)
	var transactions []*Transaction
	for rows.Next() {
		var id int64
		var owner_id int64
		var amount int64
		var description string
		rows.Scan(&id, &owner_id, &amount, &description)
		t := new(Transaction)
		t.Id = id
		t.Owner_id = owner_id
		t.Amount = amount
		t.Description = description
		transactions = append(transactions, t)
	}
	defer db.Close()
	transaction_coll := new(TransactionCollection)
	transaction_coll.Transactions = transactions
	return transaction_coll
}

func (u User) Insert() User {
	db, _ := initDB()
	stmt, err := db.Prepare("INSERT INTO users (facebook_id) VALUES ($1) RETURNING id;")
	var id int64
	err = stmt.QueryRow(u.Facebook_id).Scan(&id)
	defer db.Close()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	u.Id = id
	fmt.Printf("New user has id %d\n", u.Id)
	return u
}

type UserCollection struct {
	Users []*User
}

type Transaction struct {
	Id int64
	Owner_id int64 // foreign key to user
	Amount int64 // cents
	Description string
}

func (t Transaction) Insert() Transaction {
	db, _ := initDB()
	fmt.Printf("Inserting transaction...")
	stmt, err := db.Prepare("INSERT INTO transactions (owner_id, amount, description) VALUES ($1, $2, $3) RETURNING id;")
	var id int64
	err = stmt.QueryRow(t.Owner_id, t.Amount, t.Description).Scan(&id)
	defer db.Close()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	t.Id = id
	fmt.Printf("New Transaction inserted: %d\n", t.Id)
	return t
}

func (t Transaction) Owner() *User {
	db, _ := initDB()
	stmt, _ := db.Prepare("SELECT * FROM users WHERE id = '$1';")
	var id int64
	var facebook_id int64
	stmt.QueryRow(t.Owner_id).Scan(&id, &facebook_id)
	defer db.Close()
	user := new(User)
	user.Id = id
	user.Facebook_id = facebook_id
	return user
}

type TransactionCollection struct {
	Transactions []*Transaction
}

func NewUser(facebook_id int64) *User {
	fmt.Printf("Creating user %s\n", facebook_id)
	u := new(User)
	u.Facebook_id = facebook_id
	return u
}

func UserExists(user_id int64) bool {
	db, _ := initDB()
	stmt, _ := db.Prepare("SELECT id FROM users WHERE id = $1;")
	result, err := stmt.Exec(user_id)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	rows_affected, _ := result.RowsAffected()
	defer db.Close()
	if rows_affected > 0 {
		return true
	} 
	return false
}

func NewTransaction(owner_id int64, amount int64, description string) (*Transaction, error ){
	fmt.Printf("Creating transaction.\n")
	if !UserExists(owner_id) {
		return nil, nil
	}

	t := new(Transaction)
	t.Owner_id = owner_id
	t.Amount = amount
	t.Description = description
	return t, nil
}

// Database Interaction
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

func UserCount(db *sql.DB) int {
	var count int
	db.QueryRow("SELECT COUNT(*) FROM users;").Scan(&count)
	return count
}

func Users(db *sql.DB) (*UserCollection, error) {
	fmt.Printf("Querying for users.\n")
	rows, err := db.Query("SELECT id, facebook_id FROM users;")
	if err != nil {
		fmt.Printf("Error querying for users:\n%s\n", err)
		return nil, err
	}
	count := UserCount(db)
	fmt.Printf("Found %d users.\n", count)
	var users []*User
	for rows.Next() {
		var id int64
		var facebook_id int64
		err = rows.Scan(&id, &facebook_id)
		user := new(User)
		user.Id = id
		user.Facebook_id = facebook_id
		users = append(users, user)
	}
	user_collection := new(UserCollection)
	user_collection.Users = users
	return user_collection, nil
}

// HTTP handlers
func new_user(res http.ResponseWriter, req *http.Request) {
	db, _ := initDB()
	params := req.URL.Query()
	var facebook_id int64
	facebook_id_int, _ := strconv.Atoi(params["facebook_id"][0])
	facebook_id = int64(facebook_id_int)
	
	result := NewUser(facebook_id)
	user := result.Insert()

	user_json, _ := json.Marshal(user)
	res.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(res, string(user_json))
	defer db.Close()
}

func users(res http.ResponseWriter, req *http.Request) {
	db, err := initDB()
	if err != nil {
		fmt.Fprintf(res, "There was an error querying: %s\n", err)
	}
	users, _ := Users(db)
	users_json, _ := json.Marshal(users)
	res.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(res, string(users_json))

	defer db.Close()
}

func new_transaction(res http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	var owner_id int64
	var amount int64
	var description string

	owner_id_int, _ := strconv.Atoi(params["owner_id"][0])
	owner_id = int64(owner_id_int)

	amount_int, _ := strconv.Atoi(params["amount"][0])
	amount = int64(amount_int)

	description = params["description"][0]

	result, _ := NewTransaction(owner_id, amount, description)
	transaction := result.Insert()
	t_json, _ := json.Marshal(transaction)
	res.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(res, string(t_json))
}

// entry point
func main() {
	http.HandleFunc("/users", users)
	http.HandleFunc("/users/new", new_user)

	http.HandleFunc("/transactions/new", new_transaction)

	fmt.Printf("Listening...")
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
    if err != nil {
		panic(err)
    }
}
