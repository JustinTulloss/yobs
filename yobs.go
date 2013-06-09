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
	users := Users()
	users_json, _ := json.Marshal(users)
	res.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(res, string(users_json))
}

func transactions(res http.ResponseWriter, req *http.Request) {
	transactions := Transactions()
	t_json, _ := json.Marshal(transactions)
	res.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(res, string(t_json))
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

	fmt.Printf("Listening on localhost:%s...", os.Getenv("PORT"))
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
    if err != nil {
		panic(err)
    }
}
