package main

import (
	"fmt"
    _ "github.com/lib/pq"
    "database/sql"
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

func Transactions() *TransactionCollection {
	db, _ := initDB()

	rows, err := db.Query("SELECT id, owner_id, amount, description FROM transactions;")
	if err != nil {
		fmt.Printf("Error querying for transactions: %s\n", err)
	}
	defer db.Close()
	var transactions []*Transaction
	for rows.Next() {
		var id int64
		var owner_id int64
		var amount int64
		var description string
		rows.Scan(&id, &owner_id, &amount, &description)
		transaction := new(Transaction)
		transaction.Id = id
		transaction.Owner_id = owner_id
		transaction.Amount = amount
		transaction.Description = description
		transactions = append(transactions, transaction)
	}
	transaction_collection := new(TransactionCollection)
	transaction_collection.Transactions = transactions
	return transaction_collection
}

func Users() *UserCollection {
	fmt.Printf("Querying for users.\n")
	db, err := initDB()

	rows, err := db.Query("SELECT id, facebook_id FROM users;")
	if err != nil {
		fmt.Printf("Error querying for users:\n%s\n", err)
		return nil
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
	defer db.Close()
	user_collection := new(UserCollection)
	user_collection.Users = users
	return user_collection
}