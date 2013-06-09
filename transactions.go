package main

import (
	"fmt"
)

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

func NewTransactionByFB(facebook_id int64, amount int64, description string) (*Transaction, error) {
	u := UserFromFB(facebook_id)
	return NewTransaction(u.Id, amount, description)
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
