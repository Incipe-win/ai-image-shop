# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is an AI-powered T-shirt shop application built in Go. The project follows standard Go project layout conventions with a clean architecture pattern.

## Commands

### Build and Run
```bash
# Build the server
go build -o bin/server cmd/server/main.go

# Run the server directly
go run cmd/server/main.go

# Run with go (development)
go run ./cmd/server
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests for specific package
go test ./internal/service
```

### Development
```bash
# Format code
go fmt ./...

# Vet code for potential issues
go vet ./...

# Tidy dependencies
go mod tidy

# Download dependencies
go mod download
```

## Architecture

### Project Structure
- `cmd/server/` - Application entry point (main.go)
- `internal/` - Private application code
  - `handler/` - HTTP request handlers
  - `middleware/` - HTTP middleware components
  - `model/` - Data models and structures
  - `repository/` - Data access layer
  - `service/` - Business logic layer
- `pkg/` - Public library code (reusable packages)
- `configs/` - Configuration files

### Configuration
The application uses YAML configuration (`configs/config.yaml`) with the following main sections:
- `server.port` - Server listening port (default: ":8080")
- `database.dsn` - PostgreSQL connection string
- `jwt.secret` - JWT signing secret
- `ai.api_key` - AI service API key

### Architecture Pattern
The project follows a layered architecture:
1. **Handler Layer** - HTTP request/response handling
2. **Service Layer** - Business logic and orchestration
3. **Repository Layer** - Data persistence and retrieval
4. **Model Layer** - Data structures and domain entities

### Dependencies
- Go 1.25.0 (as specified in go.mod)
- PostgreSQL database integration
- JWT authentication
- AI service integration for t-shirt design generation

## Development Notes

- The project is in early development stage with basic structure in place
- Configuration includes placeholder values that should be updated for actual deployment
- Database schema and migrations are not yet implemented
- No existing tests or build automation detected