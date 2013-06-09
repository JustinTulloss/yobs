# Yet Another Bill Splitter

### To get set up

1. Install the prequisites: `$ gem install rails` (We use ActiveRecord for database migrations, see [more details here](http://blog.aizatto.com/2007/05/27/activerecord-migrations-without-rails/))
2. Create the database: `$ createdb yobs`
3. Run the migrations: `$ rake`
4. Seed the database: `$ psql yobs < seed.sql`
5. Start the server: `$ ./serve.sh`


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
