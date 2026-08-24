# Server Application

This example demonstrates the **in-place editing feature** of envdoc using the `-edit` flag.

## Overview

When using `-edit` mode, envdoc searches for HTML comment markers in your README file:
- `envdoc:begin` - marks the start of the generated section
- `envdoc:end` - marks the end of the generated section

All content between these markers will be replaced with newly generated documentation,
while everything else in the file remains unchanged. This allows you to maintain
environment variable documentation alongside other content in your README.

## Features

- HTTP/HTTPS server
- Database connectivity
- Configurable connection limits
- Debug logging

## Environment Variables

The following environment variables configure the server application:

<!--envdoc:begin-->

## ServerConfig

ServerConfig is the server configuration structure.
This example demonstrates using envdoc in edit mode to
maintain documentation directly within a README file.

 - `SERVER_HOST` (default: `localhost`) - Host is the server hostname or IP address to bind to.
 - `SERVER_PORT` (**required**) - Port is the server port number.
 - TLS configuration
   - `TLS_ENABLED` (default: `false`) - Enabled turns on TLS/SSL.
   - `TLS_CERT_FILE` - CertFile is the path to the TLS certificate file.
   - `TLS_KEY_FILE` - KeyFile is the path to the TLS private key file.
 - `DATABASE_URL` (**required**) - Database connection string.
 - `DEBUG` (default: `false`) - Debug enables debug logging when set to true.
 - `MAX_CONNECTIONS` (default: `100`) - MaxConnections limits the number of concurrent connections.

<!--envdoc:end-->

## Usage

```bash
# Set required environment variables
export SERVER_PORT=8080
export DATABASE_URL=postgres://localhost/mydb

# Optional configuration
export SERVER_HOST=0.0.0.0
export DEBUG=true
export MAX_CONNECTIONS=200

# TLS configuration (optional)
export TLS_ENABLED=true
export TLS_CERT_FILE=/path/to/cert.pem
export TLS_KEY_FILE=/path/to/key.pem

# Run the server
./server
```

## Development

To regenerate the environment variable documentation after modifying `config.go`:

```bash
go generate ./...
```

The `//go:generate` directive in `config.go` uses the `-edit` flag to update
only the content between the HTML comment markers, preserving all other
documentation in this README.

## How It Works

The `config.go` file contains:

```go
//go:generate go run ../../ -edit -output README.md -format markdown
```

When you run `go generate`, envdoc:
1. Parses the struct definitions and their `env` tags
2. Generates markdown documentation
3. Finds the markers in README.md
4. Replaces only the content between the markers
5. Preserves everything else in the file
