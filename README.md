# Yet Another Bill Splitter

### To get set up

1. Install the prequisites: `$ gem install rails` (We use ActiveRecord for database migrations, see [more details here](http://blog.aizatto.com/2007/05/27/activerecord-migrations-without-rails/))
2. Create the database: `$ createdb yobs`
3. Run the migrations: `$ rake`
4. Seed the database: `$ psql yobs < seed.sql`
5. Start the server: `$ ./serve.sh`


## API endpoints

###  /users
Query for a list of users.

Returns:

```
{
  "Users": [
    {
      "Facebook_id": <int>,
      "Id": <int>
    },
    {
      "Facebook_id": <int>,
      "Id": <int>
    },
    …,
  ]
}
```

### /user
Retrieve the user with the given `facebook_id` (which is required). If the user does not exist, the facebook_id will have a value of 0.

Returns:

```
{
  "Facebook_id": <ID as int>,
  "Id": <ID as int>
}
```

### /users/new
Create a new user with the given facebook ID (required).

Returns:

```
{
  "Facebook_id": <ID as int>,
  "Id": <ID as int>
}
```

### /transactions
Query for transactions. Optionally, pass a `facebook_id` to filter the results by user.

Returns:

```
{
  "Transactions": [
    {
      "Id": <int>,
      "Owner_id": <int>,
      "Amount": <int>,
      "Description": <string>
    },
	…,
  ]
}
```

### /transaction
TODO

### /transactions/new
Create a new transaction with an associated owner. Either `owner_id` or `facebook_id` is required (passing both will result in an error). You must also pass an `amount`, but `description` is **optional**.

Returns:
```
{
  "Id": <int>,
  "Owner_id": <int>,
  "Amount": <int>,
  "Description": <string>
}
```

### Schema

* Users
  * facebook_id

* StripeCustomers
  * stripe_customer_id
  * user_id

* Transactions
  * owner_id (user)
  * amount (cents)
  * description (text)

* UserTransactions
  * user_id
  * transaction_id
  * amount
