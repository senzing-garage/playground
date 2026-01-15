# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`playground` is a Go CLI application that creates an environment for exploring Senzing entity resolution. It starts concurrent gRPC and HTTP servers providing web interfaces, REST APIs, Swagger UI, and web terminal (xterm) capabilities.

## Build Commands

```bash
# Install development tools (one-time)
make dependencies-for-development

# Install Go/Python dependencies
make dependencies

# Build binary
make clean build

# Build with SQLite3 support
make build-with-libsqlite3

# Run without building
make run

# Build Docker image
make docker-build

# Run Docker container
make docker-run
```

## Testing

```bash
# Run tests
make clean setup test

# Run tests with coverage (target >80%)
make clean setup coverage

# Check coverage against threshold
make check-coverage

# Docker Compose integration tests
make docker-test
```

## Linting

```bash
# Run all linters (Go + Python)
make lint

# Individual linters
make golangci-lint    # Go linting (80+ linters via .github/linters/.golangci.yaml)
make govulncheck      # Go vulnerability check
make cspell           # Spell check

# Python linters
make pylint mypy bandit black flake8 isort

# Auto-fix many lint issues
make fix
```

## Architecture

### Multi-Server Design

The application runs two servers concurrently in [cmd/root.go](cmd/root.go):

- **HTTP Server** (port 8260): REST API, Swagger UI, Entity Search, JupyterLab, Xterm web terminal
- **gRPC Server** (port 8261): Senzing SDK services via `serve-grpc`

### Key Packages

- `cmd/` - Cobra CLI commands with Viper configuration. OS-specific context in `context_*.go`
- `httpserver/` - HTTP server implementation with route handlers for multiple services
- `rootfs/` - Docker filesystem artifacts and example Go programs

### Senzing Integration

Uses Senzing SDK components:

- `go-rest-api-service` - REST API implementation
- `serve-grpc` - gRPC server for Senzing SDK
- `demo-entity-search` - Entity search web interface
- `go-cmdhelping` - CLI helpers with standardized options

### Configuration

Environment variables follow pattern `SENZING_TOOLS_*`:

- `SENZING_TOOLS_DATABASE_URL` - Database connection
- `SENZING_TOOLS_HTTP_PORT` - HTTP port (default 8260)
- `SENZING_TOOLS_GRPC_PORT` - gRPC port (default 8261)
- `SENZING_TOOLS_LOG_LEVEL` - Logging level

## Prerequisites

Senzing C library must be installed:

- `/opt/senzing/er/lib` - Shared objects
- `/opt/senzing/er/sdk/c` - SDK headers
- `/etc/opt/senzing` - Configuration

## Senzing Guidelines

Follow instructions at <https://raw.githubusercontent.com/senzing-factory/claude/refs/tags/v1/commands/senzing.md>
