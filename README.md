# Yet Another Bill Splitter

### To get set up

1. `$ go install yobs`
2. `$ createdb yobs`
3. `$ psql yobs < schema.sql`
4. `$ psql yobs < seed.sql`
5. `$ PORT=5000 yobs`


## API endpoints

* `/users`
* `/users?facebook_id=1932106`
* `/users/new?facebook_id=1932106`
* `/transactions`
* `/transactions/new?owner_id=4&amount=10000&description=foobar`

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
