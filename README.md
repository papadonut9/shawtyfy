# Shawtyfy

![Shawtyfy logo](assets/Brand-shawtyfy.png "Shawtyfy logo")

Shawtyfy is a lightweight URL shortener service written in Go. It stores URL mappings in Redis for fast lookups and optionally syncs the data to DynamoDB for persistence.

## Prerequisites

- Go 1.21 or later
- Redis server running locally on `127.0.0.1:6379`
- (Optional) AWS credentials configured if you want DynamoDB syncing

## Installation

1. Clone the repository
   ```bash
   git clone https://github.com/papadonut9/shawtyfy.git
   cd shawtyfy
   ```
2. Download dependencies
   ```bash
   go get -d ./...
   ```
3. Start Redis
   ```bash
   redis-server
   ```
4. Run the application
   ```bash
   go run main.go
   ```
   The server listens on `http://localhost:9808`.

## Usage

Interact with the API using any HTTP client.

### Create a short URL
```bash
curl -X POST http://localhost:9808/create-short-url \
  -H "Content-Type: application/json" \
  -d '{"long_url":"https://example.com","user_id":"123"}'
```
Returns the shortened URL.

### Redirect
Visit `http://localhost:9808/<short>` to be redirected to the original address.

### Statistics and management
- `GET /get-key-count` – return number of stored URLs
- `GET /get-all-keys` – list all short codes
- `POST /delete` – remove a short code
- `POST /get-metadata` – fetch original URL and user ID

## Running tests

Run all unit tests with:
```bash
go test ./...
```

## Docker

Build the container image:
```bash
docker build -t shawtyfy .
```

Run the container:
```bash
docker run --rm -p 9808:9808 shawtyfy
```
## Contributing

Contributions are welcome! Please open an issue or pull request to propose changes or improvements.

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.
