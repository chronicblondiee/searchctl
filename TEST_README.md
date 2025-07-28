# searchctl Tests

This directory contains unit tests for the searchctl application.

## Test Structure

### Unit Tests
- **pkg/config/config_test.go** - Tests for configuration management
- **pkg/output/output_test.go** - Tests for output formatting (table, JSON, YAML)
- **pkg/client/client_test.go** - Tests for client data structures
- **internal/version/version_test.go** - Tests for version information
- **cmd/cmd_test.go** - Tests for CLI commands and help output

## Running Tests

### All Tests
```bash
make test
# or
go test ./pkg/... ./cmd/... ./internal/...
```

### Specific Package
```bash
go test ./pkg/config/
go test ./pkg/output/
go test ./pkg/client/
go test ./cmd/
go test ./internal/version/
```

### With Coverage
```bash
make test-coverage
```

### Verbose Output
```bash
go test -v ./pkg/... ./cmd/... ./internal/...
```

## Test Coverage

The tests cover:

✅ **Configuration Management**
- Default config initialization
- Context switching
- Cluster and user configuration retrieval
- Error handling for non-existent resources

✅ **Output Formatting**
- Table formatter for human-readable output
- JSON formatter for machine-readable output
- YAML formatter for configuration files
- Formatter factory function

✅ **Client Structures**
- Data structure validation
- Field mapping correctness
- Type safety

✅ **CLI Commands**
- Help command functionality
- Version command (table and JSON output)
- Dry-run operations
- Command argument parsing

✅ **Version Information**
- Version info retrieval
- Build-time variable injection
- Go version detection

## Test Philosophy

The tests are designed to be:
- **Concise** - Focus on essential functionality
- **Isolated** - Each test is independent
- **Fast** - No external dependencies in unit tests
- **Reliable** - Consistent results across environments

## Adding New Tests

When adding new functionality:
1. Create corresponding unit tests
2. Test both success and error cases
3. Use table-driven tests for multiple scenarios
4. Mock external dependencies
5. Ensure tests are deterministic

## Mocking

For testing network-dependent code:
- Use `httptest.Server` for HTTP mocking
- Create interfaces for external dependencies
- Use dependency injection for testability
