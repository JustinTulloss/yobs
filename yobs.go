package main

import (
	"fmt"
    _ "github.com/lib/pq"
    "net/http"
    "os"
	"strconv"
	"encoding/json"
)



// HTTP handlers
func new_user(res http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	var facebook_id int64
	facebook_id_int, _ := strconv.Atoi(params["facebook_id"][0])
	facebook_id = int64(facebook_id_int)
	
	result := NewUser(facebook_id)
	user := result.Insert()

	user_json, _ := json.Marshal(user)
	res.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(res, string(user_json))
}

func users(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	params := req.URL.Query()
	if len(params["facebook_id"]) > 0 {
		facebook_id_int, _ := strconv.Atoi(params["facebook_id"][0])
		// facebook id is present
		facebook_id := int64(facebook_id_int)
		user := UserFromFB(facebook_id)
		user_json, _ := json.Marshal(user)
		fmt.Fprintf(res, string(user_json))
	} else {
		// no facebook id, give back all users
		users := Users()
		users_json, _ := json.Marshal(users)
		fmt.Fprintf(res, string(users_json))
	}
}

func transactions(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	params := req.URL.Query()
	if len(params["facebook_id"]) > 0 {
		facebook_id_int, _ := strconv.Atoi(params["facebook_id"][0])
		facebook_id := int64(facebook_id_int)
		user := UserFromFB(facebook_id)
		transactions := user.Transactions()
		t_json, _ := json.Marshal(transactions)
		fmt.Fprintf(res, string(t_json))
	} else {
		transactions := Transactions()
		t_json, _ := json.Marshal(transactions)
		fmt.Fprintf(res, string(t_json))
	}
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

	http.HandleFunc("/transactions", transactions)
	http.HandleFunc("/transactions/new", new_transaction)

	fmt.Printf("Listening on localhost:%s...\n", os.Getenv("PORT"))
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
    if err != nil {
		panic(err)
    }
}
