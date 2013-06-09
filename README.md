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
