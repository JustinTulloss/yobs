package main

import (
	"fmt"
    _ "github.com/lib/pq"
    "net/http"
    "os"
	"strconv"
	"encoding/json"
)

func HasFacebookOrOwnerId(req *http.Request) (bool, string) {
	params := req.URL.Query()
	if len(params["facebook_id"]) < 1 {
		if len(params["owner_id"]) < 1 {
			return false, ""
		} else {
			return true, "owner_id"
		}
	} else {
		if !(len(params["owner_id"]) < 1) {
			// both were passed
			return false, ""
		} else {
			return true, "facebook_id"
		}
	}
	return false, ""
}

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
	res.Header().Set("Content-Type", "application/json")
	params := req.URL.Query()
	valid, key := HasFacebookOrOwnerId(req)
	if !valid {
		var errors = map[string] string {
			"error" : "Missing facebook_id or owner_id",
		}
		e_json, _ := json.Marshal(errors)
		fmt.Fprintf(res, string(e_json))
		return
	} 

	var owner_id int64
	var facebook_id int64
	var amount int64
	var description string
	
	if key == "owner_id" {
		owner_id_int, _ := strconv.Atoi(params["owner_id"][0])
		owner_id = int64(owner_id_int)
	}

	if key == "facebook_id" {
		facebook_id_int, _ := strconv.Atoi(params["facebook_id"][0])
		facebook_id = int64(facebook_id_int)
	}

	amount_int, _ := strconv.Atoi(params["amount"][0])
	amount = int64(amount_int)

	description = params["description"][0]

	var result *Transaction
	if key == "facebook_id"  {
		result, _ = NewTransactionByFB(facebook_id, amount, description)
	} else {
		result, _ = NewTransaction(owner_id, amount, description)
	}
	transaction := result.Insert()
	t_json, _ := json.Marshal(transaction)
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
