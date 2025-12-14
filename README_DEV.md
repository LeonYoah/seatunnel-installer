# SeaTunnel Enterprise Platform - Development Guide

## Project Structure

```
.
├── cmd/                          # Main applications
│   ├── installer/               # Installer API server
│   ├── agent/                   # Node agent
│   └── control-plane/           # Control plane server
├── internal/                     # Private application code
│   ├── installer/               # Installer implementation
│   ├── agent/                   # Agent implementation
│   ├── controlplane/            # Control plane implementation
│   ├── api/                     # API handlers
│   ├── service/                 # Business logic
│   ├── repository/              # Data access layer
│   └── models/                  # Data models
├── pkg/                         # Public libraries
│   ├── logger/                  # Logging utilities
│   ├── config/                  # Configuration management
│   ├── utils/                   # Common utilities
│   └── errors/                  # Error handling
├── web/                         # Frontend code (Vue3)
├── scripts/                     # Build and deployment scripts
├── docs/                        # Documentation
└── tests/                       # Test files

## Prerequisites

- Go 1.21+
- Node.js 18+ (for frontend)
- Make

## Getting Started

### 1. Install Dependencies

```bash
make deps
```

### 2. Build All Binaries

```bash
make build
```

This will create three binaries in the `bin/` directory:
- `seatunnel-installer` - Installer API server
- `seatunnel-agent` - Node agent
- `seatunnel-control-plane` - Control plane server

### 3. Run Individual Components

```bash
# Run installer server
make run-installer

# Run agent
make run-agent

# Run control plane
make run-control-plane
```

## Development Workflow

### Building

```bash
# Build all
make build

# Build specific component
make build-installer
make build-agent
make build-control-plane
```

### Testing

```bash
# Run all tests
make test

# Run tests with coverage
go test -cover ./...
```

### Cleaning

```bash
make clean
```

## Configuration

Each component can be configured via:
1. Configuration file (YAML)
2. Environment variables
3. Command-line flags

Example configuration file (`config.yaml`):

```yaml
server:
  port: 8080
  host: 0.0.0.0

database:
  type: sqlite
  database: seatunnel.db

logger:
  level: info
  output_paths:
    - stdout
    - /var/log/seatunnel/app.log
```

## Next Steps

1. Implement configuration management (Task 1.1)
2. Implement logger framework (Task 1.2)
3. Implement utility functions (Task 1.3)
4. Implement error handling (Task 1.4)
