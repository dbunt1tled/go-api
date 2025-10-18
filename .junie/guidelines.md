# Go Echo Project Development Guidelines

This document provides essential information for developers working on this Go Echo project. It covers build instructions, testing procedures, and development guidelines specific to this project.

## Build and Configuration Instructions

### Environment Setup

The project uses environment variables for configuration. Copy the `.env.example` file to `.env` and adjust the values as needed:

```bash
cp .env.example .env
```

### Build Commands

The project includes a Makefile with various commands for building and running different components:

#### API Server
```bash
# Run the API server
make run_api

# Build the API server
make build_api

# Build with optimizations
make build_opt_api
```

#### Mail Consumer
```bash
# Run the mail consumer
make run_mail

# Build the mail consumer
make build_mail

# Build with optimizations
make build_opt_mail
```

#### Centrifugo Server
```bash
# Run the Centrifugo server
make run_centrifugo

# Build the Centrifugo server
make build_centrifugo

# Build with optimizations
make build_opt_centrifugo
```

### Protocol Buffers

The project uses Protocol Buffers for gRPC communication:

```bash
# Generate Protocol Buffer code
make gen_proto

# Clean generated Protocol Buffer code
make gen_clean
```

### Database Migrations

Database migrations are managed using Goose:

```bash
# Create a new SQL migration
# Example: MIGRATION_NAME=create_table_users make migration_sql
make migration_sql

# Create a new Go migration
make migration_go

# Apply migrations
make migrate_up

# Rollback migrations
make migrate_down

# Check migration status
make migrate_status
```

## Testing Information

### Running Tests

Tests are organized by package in the `tests` directory. To run tests for a specific package:

```bash
cd tests/<package> && go test -v
```

For example, to run the Centrifugo tests:

```bash
cd tests/centrifugo && go test -v
```

### Writing Tests

Tests in this project use the [Testify](https://github.com/stretchr/testify) package for assertions and test suites.

#### Test Structure

1. Tests should be in a separate package named after the package being tested with a `_test` suffix.
2. Use table-driven tests where appropriate.
3. Use `require` for assertions that should stop the test if they fail, and `assert` for assertions that should continue.

#### Example Test

Here's an example test for the FileReader component:

```go
package reader_test

import (
	"go_echo/internal/reader"
	"go_echo/internal/reader/data"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileReader_NewFileReader(t *testing.T) {
	// Create a temporary text file for testing
	tmpFile, err := os.CreateTemp("", "test_*.txt")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	
	// Write some test data
	_, err = tmpFile.WriteString("column1,column2,column3\nvalue1,value2,value3\nvalue4,value5,value6")
	require.NoError(t, err)
	require.NoError(t, tmpFile.Close())
	
	// Test cases
	tests := []struct {
		name        string
		options     data.FileOptions
		expectError bool
	}{
		{
			name: "Valid TXT file",
			options: data.FileOptions{
				FileName: tmpFile.Name(),
			},
			expectError: false,
		},
		{
			name: "Unsupported file extension",
			options: data.FileOptions{
				FileName: "test.unsupported",
			},
			expectError: true,
		},
	}
	
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a new file reader
			fr, err := reader.NewFileReader(tc.options)
			
			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, fr)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, fr)
				
				// Test reading lines
				lines, errCh := fr.ReadLine()
				
				// Read the first line
				line, ok := <-lines
				assert.True(t, ok)
				assert.Contains(t, line, "column1")
				
				// Check for errors
				select {
				case err := <-errCh:
					assert.NoError(t, err)
				default:
					// No error, which is expected
				}
			}
		})
	}
}
```

#### Test Suites

For more complex tests, use test suites. The project includes a suite implementation in `tests/centrifugo/suite/suite.go` that can be used as a reference.

### Vulnerability Checking

The project includes a command to check for vulnerabilities:

```bash
# Install govulncheck
make install_govulncheck

# Check for vulnerabilities
make check_vulnerabilities
```

## Code Style and Development Guidelines

### Linting

The project uses [golangci-lint](https://golangci-lint.run/) for linting with a comprehensive set of linters enabled. The configuration is in `.golangci.yml`.

To run the linter:

```bash
golangci-lint run
```

### Code Style Guidelines

1. **Error Handling**: Always check errors and return them to the caller when appropriate.
2. **Function Length**: Keep functions under 100 lines and 50 statements.
3. **Cognitive Complexity**: Keep cognitive complexity under 20.
4. **Magic Numbers**: Avoid magic numbers; use constants or named variables.
5. **Comments**: End comments with a period.
6. **Imports**: Use goimports to organize imports.
7. **Context**: Use context in functions that perform I/O operations.

### Project Structure

The project follows a modular structure:

- `app/`: Application-specific code organized by domain
- `cmd/`: Entry points for different executables
- `internal/`: Internal packages not meant to be imported by other projects
- `resources/`: Static resources like templates and localization files
- `tests/`: Test files organized by package

### Logging

The project uses a custom logger defined in `internal/lib/logger/logger.go`. Use this logger instead of the standard library's log package.

### Configuration

Configuration is loaded from environment variables using the `internal/config/env/env.go` package. Add new configuration options to this package.

### Internationalization

The project supports multiple languages with localization files in the `resources/` directory. Use the `internal/config/locale/locale.go` package for internationalization.

### Authentication

Authentication is handled by the `app/auth/` package. Use the authentication middleware in `internal/router/middlewares/auth.bearer.go` for protected routes.

### Real-time Communication

The project uses Centrifugo for real-time communication. See the `app/centrifugo/` package for implementation details.
