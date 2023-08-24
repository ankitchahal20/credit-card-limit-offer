# Credit Card Offer Limit Application 

This repository contains the source code for credit card offer limiy application built using Golang. The system is responsible for creating an credit account. Once the credit account is created, one can create, accept, reject and view the credit limit offer for that account.

## Prerequisites

Before running the Credit Card Limit Offer, make sure you have the following prerequisites installed on your system:

- Go programming language (go1.20.4)
- PostgreSQL(14.8)

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/ankitchahal20/credit-card-limit-offer.git
   ```

2. Navigate to the project directory:

   ```bash
   cd credit-card-limit-offer
   ```

3. Install the required dependencies:

   ```bash
   go mod tidy
   ```

4. DB setup
    ```
    Use the scripts inside sql-scripts directory to create the tables in your db.
    ```
5. Defaults.toml
Add the values to defaults.toml and execute `go run main.go` from the cmd directory.

## APIs
There are five API's which this repo currently supports.

Create Account API
```
curl -i -k -X POST \
   http://localhost:8080/v1/create_account \
  -H "transaction-id: 288a59c1-b826-42f7-a3cd-bf2911a5c351" \
  -H "content-type: application/json" \
  -d '{
  "account_limit": 1000,
  "per_transaction_limit": 1000,
  "last_account_limit": 1000,
  "last_per_transaction_limit": 1000
}'
```

Get Account API

```
curl -i -k -X POST \
  http://localhost:8080/v1/login \
  -H "transaction-id: 288a59c1-b826-42f7-a3cd-bf2911a5c351" \
  -H "content-type: application/json" \
  -d '{
"email": "abcd11@gmail.com",
"password": "12345"
}'
```

Create Limit Offer API

```
curl -i -k -X GET \
  http://localhost:8080/v1/create_limit_offer \
  -H "transaction-id: 288a59c1-b826-42f7-a3cd-bf2911a5c351" \
  -H "content-type: application/json" \
  -d '{
  "account_id": "2b4e1e64-624f-4a4e-9911-e0b13f526e10",
  "limit_type": "PER_TRANSACTION_LIMIT",
  "new_limit": 20000,
  "offer_activation_time": "2023-08-24T02:17:00.017507+05:30",
  "offer_expiry_time": "2023-08-24T02:24:00.017507+05:30"
}
'
```

List Active Limit Offer API

```
curl -i -k -X GET \
  http://localhost:8080/v1/list_active_limit_offers \
  -H "transaction-id: 288a59c1-b826-42f7-a3cd-bf2911a5c351" \
  -H "content-type: application/json" \
  -d '{
  "account_id": "2b4e1e64-624f-4a4e-9911-e0b13f526e10",
  "active_date": "2023-08-24T02:24:00.017507+05:30"
}'
```

Update Limit Offer Status API

```
curl -i -k -X PATCH \
  http://localhost:8080/v1/update_limit_offer_status \
  -H "transaction-id: 288a59c1-b826-42f7-a3cd-bf2911a5c351" \
  -H "content-type: application/json" \
  -d '{
  "limit_offer_id": "abfe7bda-d59d-49ca-b2d0-39e3ecb12fb5",
  "status": "ACCEPTED"
  }'
```

## Project Structure

The project follows a standard Go project structure:

- `config/`: Configuration file for the application.
- `internal/`: Contains the internal packages and modules of the application.
  - `config/`: Global configuration which can be used anywhere in the application.
  - `constants/`: Contains constant values used throughout the application.
  - `db/`: Contains the database package for interacting with PostgreSQL.
  - `middleware`: Contains the logic to validate the incoming request
  - `models/`: Contains the data models used in the application.
  - `limitoffererror`: Defines the errors in the application
  - `service/`: Contains the business logic and services of the application.
  - `server/`: Contains the server logic of the application.
  - `utils/`: Contains utility functions and helpers.
- `cmd/`:  Contains command you want to build.
    - `main.go`: Main entry point of the application.
- `README.md`: README.md contains the description for the notes-taking-application.

## Contributing

Contributions to the Credit Card Limit Offer are welcome. If you find any issues or have suggestions for improvement, feel free to open an issue or submit a pull request.

## License

The Credit Card Limit Offer is open-source and released under the [MIT License](LICENSE).

## Contact

For any inquiries or questions, please contact:

- Ankit Chahal
- ankitchahal20@gmail.com

Feel free to reach out with any feedback or concerns.
