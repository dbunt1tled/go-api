# GoAPI (Echo+MySQL)

[![Go Version](https://img.shields.io/github/go-mod/go-version/dbunt1tled/go-api)](https://golang.org/)
[![Go Reference](https://pkg.go.dev/badge/github.com/dbunt1tled/go-api.svg)](https://pkg.go.dev/github.com/dbunt1tled/go-api)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/dbunt1tled/go-api)](https://goreportcard.com/report/github.com/dbunt1tled/go-api)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/dbunt1tled/go-api)](https://github.com/dbunt1tled/go-api/releases)
[![Build Status](https://github.com/dbunt1tled/go-api/actions/workflows/release.yml/badge.svg)](https://github.com/dbunt1tled/go-api/actions/workflows/release.yml)

A comprehensive boilerplate for REST API with real-time capabilities, built with Go Echo framework and MySQL database.

![GoAPI](assets/images/goapi.jpeg)

## Overview

GoAPI is a production-ready boilerplate for building REST APIs with real-time messaging capabilities. It provides a solid foundation with authentication, authorization, file processing, email services, and WebSocket communication through Centrifugo.

### Features

- **Authentication & Authorization**: JWT-based authentication with RBAC (Role-Based Access Control)
- **User Management**: User registration, profile management, and email confirmation
- **Real-time Communication**: WebSocket support via Centrifugo gRPC server
- **File Processing**: Support for XLSX, CSV, and TXT file reading/writing
- **Email Services**: Email templates and asynchronous mail sending via RabbitMQ
- **Database Migrations**: Goose-based database migration system
- **Internationalization**: Multi-language support (English/Russian)
- **Comprehensive Logging**: Structured logging with pretty formatting
- **Security**: Input sanitization, validation, and CORS support
- **Monitoring**: Optional profiling and vulnerability checking

## Technology Stack

### Core Technologies
- **Language**: Go 1.24.1
- **Web Framework**: Echo v4.13.4
- **Database**: MySQL with go-sql-driver/mysql
- **Cache**: Redis v9.14.0
- **Message Queue**: RabbitMQ with go-rabbitmq
- **Real-time**: Centrifugo with gRPC

### Key Dependencies
- **Authentication**: JWT with golang-jwt/jwt/v5
- **Validation**: go-playground/validator/v10
- **File Processing**: excelize/v2, xlsxreader
- **Email**: go-mail
- **Testing**: Testify
- **Migrations**: Goose
- **Internationalization**: go-i18n/v2

## Requirements

- **Go**: 1.24.1 or higher
- **MySQL**: 5.7+ or 8.0+
- **Redis**: 6.0+
- **RabbitMQ**: 3.8+
- **Centrifugo**: v5+ (optional, for real-time features)

## Installation & Setup

### 1. Clone the Repository

```bash
git clone https://github.com/dbunt1tled/go-api.git
cd go-api
```

### 2. Environment Configuration

Copy the environment example file and configure your settings:

```bash
cp .env.example .env
```

Edit the `.env` file with your configuration:

```env
# Application
ENV="dev"
APP_NAME="GoAPI"
APP_URL="https://127.0.0.1:8082"
APP_LOCAL="en"

# Database
DATABASE_DSN="user:password@tcp(localhost:3306)/database_name?charset=utf8mb4&parseTime=True&loc=Local"

# HTTP Server
HTTP_SERVER_ADDRESS="localhost:8082"
HTTP_SERVER_TIMEOUT=4s
HTTP_SERVER_IDLE_TIMEOUT=60s

# Redis
REDIS_HOST=127.0.0.1
REDIS_PORT=6379
REDIS_DB=0

# JWT Authentication
JWT_PRIVATE_KEY="path/to/private.key"
JWT_PUBLIC_KEY="path/to/public.key"
JWT_TOKEN_ALGORITHM="ES512"

# Email Configuration
MAIL_HOST=smtp.gmail.com
MAIL_PORT=587
MAIL_FROM_ADDRESS="your-email@example.com"
MAIL_USERNAME="your-username"
MAIL_PASSWORD="your-password"

# RabbitMQ
RABBITMQ_HOST=localhost
RABBITMQ_PORT=5672
RABBITMQ_USER="guest"
RABBITMQ_PASSWORD="guest"

# Centrifugo (optional)
SERVER_CENTRIFUGO_URL="localhost:5001"
API_CENTRIFUGO_URL=http://127.0.0.1:8000
API_CENTRIFUGO_KEY="your-centrifugo-key"
```

### 3. Install Dependencies

```bash
go mod download
```

### 4. Database Setup

Run database migrations:

```bash
make migrate_up
```

## Available Scripts

The project includes a comprehensive Makefile with the following commands:

### Application Services

```bash
# Run the main REST API server
make run_api

# Run the mail consumer service
make run_mail

# Run the Centrifugo gRPC server
make run_centrifugo
```

### Building

```bash
# Build API server
make build_api

# Build mail consumer
make build_mail

# Build Centrifugo server
make build_centrifugo

# Build with optimizations (smaller binaries)
make build_opt_api
make build_opt_mail
make build_opt_centrifugo
```

### Database Migrations

```bash
# Create new SQL migration
MIGRATION_NAME=create_users_table make migration_sql

# Create new Go migration
MIGRATION_NAME=seed_initial_data make migration_go

# Apply migrations
make migrate_up

# Rollback migrations
make migrate_down

# Check migration status
make migrate_status
```

### Protocol Buffers

```bash
# Generate Protocol Buffer code
make gen_proto

# Clean generated Protocol Buffer code
make gen_clean
```

### Security & Quality

```bash
# Install vulnerability checker
make install_govulncheck

# Check for vulnerabilities
make check_vulnerabilities
```

## Entry Points

The project has multiple entry points for different services:

### 1. REST API Server (`cmd/api/api.go`)
Main HTTP server providing REST API endpoints with:
- User authentication and management
- File upload/download functionality
- Health checks and monitoring
- Static file serving

**Default Address**: `localhost:8082`

### 2. Mail Consumer (`cmd/mailconsumer/mail.go`)
Background service for processing email queue messages:
- User registration confirmation emails
- Password reset notifications
- System notifications

### 3. Centrifugo gRPC Server (`cmd/centrifugo/centrifugo_server.go`)
Real-time communication server providing:
- WebSocket connections
- Real-time notifications
- Live data updates

**Default Address**: `localhost:5001`

### 4. Additional Utilities
- `cmd/pusher/pusher.go` - Message pushing utility
- `cmd/read/read.go` - File reading utility
- `cmd/xlsx/xslx.go` - Excel processing utility

## Environment Variables

### Application Settings
| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `ENV` | Environment (dev/prod) | `dev` | No |
| `APP_NAME` | Application name | `GoAPI` | No |
| `APP_URL` | Application URL | - | Yes |
| `APP_LOCAL` | Default language | `en` | No |

### Database & Cache
| Variable | Description | Required |
|----------|-------------|----------|
| `DATABASE_DSN` | MySQL connection string | Yes |
| `REDIS_HOST` | Redis host | Yes |
| `REDIS_PORT` | Redis port | Yes |
| `REDIS_DB` | Redis database number | No |

### HTTP Server
| Variable | Description | Default |
|----------|-------------|---------|
| `HTTP_SERVER_ADDRESS` | Server bind address | `localhost:8082` |
| `HTTP_SERVER_TIMEOUT` | Request timeout | `4s` |
| `HTTP_SERVER_IDLE_TIMEOUT` | Idle timeout | `60s` |
| `TLS_CERT` | TLS certificate path | - |
| `TLS_KEY` | TLS private key path | - |

### Authentication
| Variable | Description | Required |
|----------|-------------|----------|
| `JWT_PRIVATE_KEY` | JWT private key | Yes |
| `JWT_PUBLIC_KEY` | JWT public key | Yes |
| `JWT_TOKEN_ALGORITHM` | JWT algorithm | No |
| `SYSTEM_API_KEY` | System API key | Yes |

### External Services
| Variable | Description | Required |
|----------|-------------|----------|
| `MAIL_HOST` | SMTP host | Yes |
| `MAIL_PORT` | SMTP port | Yes |
| `MAIL_USERNAME` | SMTP username | Yes |
| `MAIL_PASSWORD` | SMTP password | Yes |
| `RABBITMQ_HOST` | RabbitMQ host | Yes |
| `RABBITMQ_PORT` | RabbitMQ port | Yes |
| `RABBITMQ_USER` | RabbitMQ username | Yes |
| `RABBITMQ_PASSWORD` | RabbitMQ password | Yes |

## Testing

The project uses [Testify](https://github.com/stretchr/testify) for testing with a comprehensive test suite.

### Running Tests

```bash
# Run tests for specific package
cd tests/reader && go test -v

# Run Centrifugo tests
cd tests/centrifugo && go test -v

# Run all tests
go test ./tests/...
```

### Test Organization

Tests are organized in the `tests/` directory by package:
- `tests/reader/` - File reading/processing tests
- `tests/centrifugo/` - Real-time communication tests

Tests follow these conventions:
- Use table-driven tests where appropriate
- Use `require` for critical assertions
- Use `assert` for non-critical assertions
- Include test suites for complex scenarios

## Project Structure

```
├── app/                          # Application logic by domain
│   ├── auth/                     # Authentication handlers & services
│   ├── centrifugo/              # Centrifugo gRPC server implementation
│   ├── general/                 # General purpose handlers
│   ├── jobs/                    # Background job handlers
│   ├── user/                    # User management
│   └── usernotification/        # User notification system
├── assets/                      # Static assets (images, styles)
├── bin/                         # Binary files and certificates
├── cmd/                         # Application entry points
│   ├── api/                     # REST API server
│   ├── centrifugo/             # gRPC server for real-time
│   ├── mailconsumer/           # Email processing service
│   └── [other utilities]/      # Additional command-line tools
├── internal/                    # Internal packages
│   ├── cache/                   # Redis cache implementation
│   ├── config/                  # Configuration management
│   ├── database/migrations/     # Database migration files
│   ├── dto/                     # Data Transfer Objects
│   ├── lib/                     # Internal libraries
│   ├── reader/                  # File reading implementations
│   ├── router/                  # HTTP routing and middleware
│   ├── storage/                 # Database connection management
│   ├── util/                    # Utility functions
│   └── writer/                  # File writing implementations
├── proto/                       # Protocol Buffer definitions
├── resources/                   # Templates and localization
│   ├── en/                      # English translations
│   ├── ru/                      # Russian translations
│   └── templates/               # HTML templates
└── tests/                       # Test files organized by package
```

## Development Guidelines

### Code Style
- Follow standard Go conventions
- Use golangci-lint for code quality
- Keep functions under 100 lines
- Maintain cognitive complexity under 20
- Always handle errors appropriately

### Git Hooks
The project includes automated workflows:
- **Release Workflow**: Automated releases with GoReleaser
- **Junie Workflow**: AI-assisted development workflow

### Linting
Run linting with:
```bash
golangci-lint run
```

## API Documentation

<!-- TODO: Add API documentation link or generate with Swagger -->
API documentation is available at:
- **Development**: `http://localhost:8082/docs` (TODO: Implement Swagger)
- **Swagger Spec**: `./docs/swagger.yaml` (TODO: Generate)

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Commit your changes: `git commit -m 'Add amazing feature'`
4. Push to the branch: `git push origin feature/amazing-feature`
5. Open a Pull Request

### TODO Items

- [ ] Add Swagger/OpenAPI documentation
- [ ] Implement health check endpoints
- [ ] Add Docker and Docker Compose configuration
- [ ] Create deployment guides for different platforms
- [ ] Add performance benchmarks
- [ ] Implement rate limiting
- [ ] Add more comprehensive integration tests
- [ ] Create development environment setup scripts

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- **Issues**: [GitHub Issues](https://github.com/dbunt1tled/go-api/issues)
- **Discussions**: [GitHub Discussions](https://github.com/dbunt1tled/go-api/discussions)

## Acknowledgments

- [Echo](https://echo.labstack.com/) - High performance, extensible web framework
- [Centrifugo](https://centrifugal.dev/) - Real-time messaging server
- [Testify](https://github.com/stretchr/testify) - Testing toolkit
- [golangci-lint](https://golangci-lint.run/) - Go linter

---

**Last Updated**: 2025-01-18
