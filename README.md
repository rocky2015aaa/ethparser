# Ethparser

Ethereum blockchain parser that will allow to query transactions for
subscribed addresses.

## Endpoints
* `POST /subscribe` Allows a user to subscribe to an Ethereum address to track transactions 
* `GET /transactions` Retrieves a list of transactions for an Ethereum address
* `GET /openapi.html` api document

## How to Run

- Generate API doc: `make api-docs`.
- Run unit tests: `make test`
- Run the server: `make run`