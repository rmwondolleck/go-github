# go-github

Home Lab API Service - A RESTful API for managing home lab devices and services.

## Features

- RESTful API endpoints
- Health check monitoring
- Request ID tracking
- Structured logging with slog
- Graceful shutdown
- Comprehensive API documentation with Swagger/OpenAPI

## API Documentation

The API is documented using OpenAPI/Swagger specification. Once the server is running, you can access:

- **Swagger UI**: [http://localhost:8080/api/docs/index.html](http://localhost:8080/api/docs/index.html)
  - Interactive API documentation
  - "Try it out" functionality to test endpoints directly
  - Full endpoint specifications

- **OpenAPI Specification**: [http://localhost:8080/api/docs/doc.json](http://localhost:8080/api/docs/doc.json)
  - Download the OpenAPI spec in JSON format
  - Can be imported into API testing tools (Postman, Insomnia, etc.)

## Getting Started

### Prerequisites

- Go 1.25.0 or later
- Make (optional, for using Makefile commands)

### Installation

```bash
# Clone the repository
git clone https://github.com/rmwondolleck/go-github.git
cd go-github

# Install dependencies
go mod download
```

### Running the Application

```bash
# Using Make
make run

# Or directly with Go
go run ./cmd/api

# Or build and run
make build
./bin/homelab-api
```

The server will start on port 8080 by default. You can change the port by setting the `PORT` environment variable:

```bash
PORT=3000 ./bin/homelab-api
```

## Development

### Building

```bash
make build
```

### Running Tests

```bash
make test
```

### Generating Swagger Documentation

```bash
make swagger
```

### Linting

```bash
make lint
```

## API Endpoints

### Health Check

```
GET /health
```

Returns the health status of the API.

### API Root

```
GET /api/v1
```

Returns API version information.

## Architecture

- **cmd/api**: Application entry point
- **internal/server**: HTTP server setup and routing
- **internal/middleware**: Middleware components (logging, recovery, request ID)
- **internal/handlers**: HTTP request handlers
- **internal/models**: Data models
- **api**: Generated Swagger/OpenAPI documentation

## License

Apache 2.0