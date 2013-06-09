package main

import (
	"fmt"
    "database/sql"
)

type User struct {
	Facebook_id int64
	Id int64
}

func (u User) Transactions() *TransactionCollection {
	db, _ := initDB()
	stmt, err := db.Prepare("SELECT id, owner_id, amount, description FROM transactions WHERE owner_id = $1;")
	if err != nil {
		fmt.Printf("%s\n",err)
	}
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

func UserCount(db *sql.DB) int {
	var count int
	db.QueryRow("SELECT COUNT(*) FROM users;").Scan(&count)
	return count
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

func UserFromFB(facebook_id int64) *User {
	db, _ := initDB()
	stmt, _ := db.Prepare("SELECT id FROM users WHERE facebook_id = $1")
	var id int64
	stmt.QueryRow(facebook_id).Scan(&id)
	user := new(User)
	user.Id = id
	user.Facebook_id = facebook_id
	defer db.Close()
	return user
}