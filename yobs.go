package main

import (
	"fmt"
    _ "github.com/lib/pq"
    "net/http"
    "os"
	"strconv"
	"encoding/json"
	"log"
)

func ToJson(v interface{}) []byte{
	result, _ := json.MarshalIndent(v, "", "  ")
	return result
}

func LogRequest(req *http.Request) {
	log.Printf("%s %s\n", req.Method, req.URL)
}

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
	LogRequest(req)
	params := req.URL.Query()
	var facebook_id int64
	facebook_id_int, _ := strconv.Atoi(params["facebook_id"][0])
	facebook_id = int64(facebook_id_int)
	
	result := NewUser(facebook_id)
	user := result.Insert()

	user_json := ToJson(user)
	res.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(res, string(user_json))
}

func user(res http.ResponseWriter, req *http.Request) {
	LogRequest(req)
	params := req.URL.Query()
	if len(params["facebook_id"]) > 0 {
		facebook_id_int, _ := strconv.Atoi(params["facebook_id"][0])
		// facebook id is present
		facebook_id := int64(facebook_id_int)
		user := UserFromFB(facebook_id)
		user_json := ToJson(user)
		fmt.Fprintf(res, string(user_json))
	} else if len(params["id"]) > 0 {
		owner_id_int, _ := strconv.Atoi(params["id"][0])
		owner_id := int64(owner_id_int)
		user := UserFromID(owner_id)
		user_json := ToJson(user)
		fmt.Fprintf(res, string(user_json))
	} else {
		var errors = map[string] string {
			"error" : "Missing facebook_id or id",
		}
		e_json := ToJson(errors)
		fmt.Fprintf(res, string(e_json))
	}
}

func users(res http.ResponseWriter, req *http.Request) {
	LogRequest(req)
	res.Header().Set("Content-Type", "application/json")
	users := Users()
	users_json := ToJson(users)
	fmt.Fprintf(res, string(users_json))
}

func transaction(res http.ResponseWriter, req *http.Request) {
	LogRequest(req)
	res.Header().Set("Content-Type", "application/json")
	params := req.URL.Query()
	if len(params["id"]) > 0 {
		id_int, _ := strconv.Atoi(params["id"][0])
		id := int64(id_int)
		transaction := TransactionFromID(id)
		t_json := ToJson(transaction)
		fmt.Fprintf(res, string(t_json))
	} else {
		// missing ID parameter
		var errors = map[string] string {
			"error" : "Missing id",
		}
		e_json := ToJson(errors)
		fmt.Fprintf(res, string(e_json))
	}
}

func transactions(res http.ResponseWriter, req *http.Request) {
	LogRequest(req)
	res.Header().Set("Content-Type", "application/json")
	params := req.URL.Query()
	if len(params["facebook_id"]) > 0 {
		facebook_id_int, _ := strconv.Atoi(params["facebook_id"][0])
		facebook_id := int64(facebook_id_int)
		user := UserFromFB(facebook_id)
		transactions := user.Transactions()
		t_json := ToJson(transactions)
		fmt.Fprintf(res, string(t_json))
	} else {
		transactions := Transactions()
		t_json := ToJson(transactions)
		fmt.Fprintf(res, string(t_json))
	}
}

func new_transaction(res http.ResponseWriter, req *http.Request) {
	LogRequest(req)
	res.Header().Set("Content-Type", "application/json")
	params := req.URL.Query()
	valid, key := HasFacebookOrOwnerId(req)
	if !valid {
		var errors = map[string] string {
			"error" : "Missing facebook_id or owner_id",
		}
		e_json := ToJson(errors)
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
	t_json := ToJson(transaction)
	fmt.Fprintf(res, string(t_json))
}

// entry point
func main() {
	http.HandleFunc("/users", users)
	http.HandleFunc("/user", user)
	http.HandleFunc("/users/new", new_user)

	http.HandleFunc("/transactions", transactions)
	http.HandleFunc("/transaction", transaction)
	http.HandleFunc("/transactions/new", new_transaction)

	log.Printf("Listening on localhost:%s...\n", os.Getenv("PORT"))
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
    if err != nil {
		panic(err)
    }
}
