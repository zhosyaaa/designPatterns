# Software design patterns projects

This project is a demonstration of how to use design patterns in Golang. The project is a simple banking application that allows users to create accounts, make payments, and view transactions.

---
### The project is structured using the following design patterns:

* Singleton: The project utilizes the Singleton pattern for the database connection by ensuring only one instance of the database is created using sync.Once.
* Strategy: The Strategy pattern is implemented to handle different payment methods (`KaspiPayment` and `PayPalPayment`). The `Purchase` struct acts as a context to switch between payment strategies.
* Factory: The CurrencyService and TransactionFactory classes use the factory pattern to create instances of different currency services and transaction types.
* Adapter: The CurrencyClientAdapter class uses the adapter pattern to adapt the CurrencyClient interface to the CurrencyService interface.
* Decorator: The BalanceCheckingDecorator and CurrencyConversionDecorator classes use the decorator pattern to add additional functionality to the Payment interface.
Project structure
* Observer: Users can subscribe or unsubscribe to specific currency updates and receive notifications via email upon currency changes.


## Installation and Setup

### PostgreSQL Database

To utilize the PostgreSQL database in the application, Docker can be used for installation and setup:

1. Launch the PostgreSQL 15 container using Docker:

    ```bash
    docker run --rm -d --name postgres15 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=1079 -p 5432:5432 postgres:15
    ```

2. Create the `bank` database:

    ```bash
    docker exec -it postgres15 createdb --username=postgres --owner=postgres bank
    ```

3. (Optional) To drop the `bank` database, execute:

    ```bash
    docker exec -it postgres15 dropdb -U postgres bank
    ```

### Running the Application

1. Install all necessary dependencies:

    ```bash
    go get -d ./...
    ```

2. Run the application:

    ```bash
    go run .
    ```

After completing these steps, your application will be running and accessible at http://localhost:8080.

---
## Endpoints

### Account Endpoints

#### Create Account

- **Endpoint:** `POST /accounts`
- **Description:** Creates a new account for a user.
- **Request Body:**
    ```json
    {
        "username": "example_user",
        "email": "user@example.com",
        "initialBalance": 1000.00
    }
    ```
- **Response:** `201 Created`
- **Example Response Body:**
    ```json
    {
        "accountId": "1234567890",
        "username": "example_user",
        "email": "user@example.com",
        "balance": 1000.00
    }
    ```

#### Get Account Details

- **Endpoint:** `GET /accounts/{accountId}`
- **Description:** Retrieves account details based on the account ID.
- **Response:** `200 OK`
- **Example Response Body:**
    ```json
    {
        "accountId": "1234567890",
        "username": "example_user",
        "email": "user@example.com",
        "balance": 1000.00
    }
    ```

### Transaction Endpoints

#### Make Payment

- **Endpoint:** `POST /transactions`
- **Description:** Initiates a payment from one account to another.
- **Request Body:**
    ```json
    {
        "fromAccountId": "1234567890",
        "toAccountId": "0987654321",
        "amount": 100.00,
        "description": "Payment for services"
    }
    ```
- **Response:** `200 OK`
- **Example Response Body:**
    ```json
    {
        "transactionId": "9876543210",
        "fromAccountId": "1234567890",
        "toAccountId": "0987654321",
        "amount": 100.00,
        "description": "Payment for services",
        "timestamp": "2023-11-21T12:00:00Z"
    }
    ```

#### Get Transaction Details

- **Endpoint:** `GET /transactions/{transactionId}`
- **Description:** Retrieves transaction details based on the transaction ID.
- **Response:** `200 OK`
- **Example Response Body:**
    ```json
    {
        "transactionId": "9876543210",
        "fromAccountId": "1234567890",
        "toAccountId": "0987654321",
        "amount": 100.00,
        "description": "Payment for services",
        "timestamp": "2023-11-21T12:00:00Z"
    }
    ```


