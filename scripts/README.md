# Test Scripts

Automated testing for searchctl functionality with emoji-free logging and shared utilities for better CI/CD compatibility.

## Quick Start

```bash
# Setup test environment
./scripts/start-test-env.sh

# Run core tests (recommended for daily use)
./scripts/integration-test.sh        # Core functionality (~30s)
./scripts/test-conditions.sh        # Conditions validation (~30s)

# Run performance tests (for regression testing)
./scripts/test-performance.sh       # Performance benchmarks (~1m)

# Run real operations (use with caution - creates actual data)
./scripts/test-rollover-real.sh     # Real rollover operations (~2m)

# Cleanup
./scripts/stop-test-env.sh
```

### Quick Validation
```bash
# Just verify everything is working
./scripts/check-status.sh           # Check container health (~5s)
./scripts/test-config.sh           # Verify configuration (~10s)
```

### CI/CD Usage
```bash
# Automated testing pipeline
make test-integration               # Safe tests only
make test-all                      # Full test suite
```

## Scripts Overview

| Script | Purpose | Safety | Duration | Usage |
|--------|---------|---------|----------|-------|
| `integration-test.sh` | Core functionality validation | ✅ Safe (dry-run) | ~30s | Daily CI |
| `test-rollover-real.sh` | Real rollover operations | ⚠️ Creates data | ~2m | Pre-release testing |
| `test-performance.sh` | Performance benchmarking | ✅ Safe (dry-run) | ~1m | Performance regression |
| `test-conditions.sh` | Conditions file validation | ✅ Safe (dry-run) | ~30s | Config testing |
| `test-config.sh` | Configuration testing | ✅ Safe (read-only) | ~10s | Config validation |
| `start-test-env.sh` | Environment setup | ✅ Safe | ~30s | Setup |
| `stop-test-env.sh` | Environment cleanup | ✅ Safe | ~10s | Cleanup |
| `check-status.sh` | Health verification | ✅ Safe | ~5s | Debugging |

## Architecture

### Common Utilities (`common.sh`)
- **Consistent Logging**: Color-coded output with `[INFO]`, `[SUCCESS]`, `[ERROR]` prefixes
- **Environment Management**: Automated setup, health checks, cleanup
- **Test Execution**: Standardized command execution with timing
- **Performance Testing**: Benchmarking utilities for load testing
- **Extensibility**: Reusable functions for future test development

### Test Categories

#### Core Features
- **Cluster Operations**: health, info, connectivity
- **Resource Management**: indices, nodes, data streams
- **CRUD Operations**: create, get, delete resources

#### Rollover Features  
- **Conditions**: age, docs, size, primary shard limits
- **Output Formats**: table, json, yaml
- **Advanced Options**: lazy rollover, conditions files
- **Error Handling**: missing args, invalid formats

#### Data Stream Features
- **Lifecycle**: create, list, delete data streams
- **Integration**: with rollover operations
- **Patterns**: wildcard matching, filtering

### Test Output Format

All scripts use consistent emoji-free logging:

```bash
[INFO] Checking test environment...
[SUCCESS] Elasticsearch is ready
[SUCCESS] OpenSearch is ready
[TEST] Testing elasticsearch...
[EXEC] Running: ./bin/searchctl --context elasticsearch cluster health
[EXEC] Command succeeded
[SUCCESS] Integration tests completed successfully!
```

**Log Levels:**
- `[INFO]` - General information and progress updates
- `[SUCCESS]` - Successful operations and completions
- `[ERROR]` - Failures and error conditions  
- `[TEST]` - Test operation descriptions
- `[EXEC]` - Command execution details
- `[BUILD]` - Build and compilation messages
- `[TIMING]` - Performance timing information

## Configuration

### Environment Variables
```bash
export SEARCHCTL_CONFIG="examples/test-config.yaml"  # Config file path
export SEARCHCTL_CONTEXT="elasticsearch"            # Force context
export TEST_TIMEOUT="60s"                          # Test timeout
```

### Test Environment
- **Elasticsearch**: `localhost:9200`
- **OpenSearch**: `localhost:9201` 
- **Security**: Disabled for testing
- **Resources**: 512MB RAM per service

## Adding New Tests

### 1. Create Test Script
Use the common utilities for consistent functionality:

```bash
#!/bin/bash
set -e

# Source common utilities
source "$(dirname "$0")/common.sh"

# Setup environment  
setup_test_environment
check_environment

# Your tests here
test_new_feature() {
    log_info "Testing new feature..."
    test_command "./bin/searchctl new-command --dry-run"
    log_success "New feature test completed"
}

# Run tests
test_new_feature
```

### 2. Available Utilities

#### Core Functions
- `setup_test_environment()` - Build searchctl and configure test environment
- `check_environment()` - Verify Elasticsearch and OpenSearch are running
- `cleanup_test_data()` - Clean up temporary files and test data

#### Logging Functions  
- `log_info()` - Blue `[INFO]` messages for general information
- `log_success()` - Green `[SUCCESS]` messages for completed operations
- `log_error()` - Red `[ERROR]` messages for failures
- `log_test()` - Yellow `[TEST]` messages for test operations

#### Test Execution
- `test_command()` - Execute commands with error handling and logging
- `time_command()` - Execute commands with performance timing
- `benchmark_command()` - Run performance benchmarks with iterations
- `test_both_engines()` - Test functionality against both ES and OpenSearch

#### Validation Utilities
- `validate_json()` - Verify JSON output is well-formed
- `validate_yaml()` - Verify YAML output is well-formed
- `wait_for_service()` - Wait for services to become available

### 3. Update Integration Script
Add new test calls to `integration-test.sh`:
```bash
log_info "Testing new feature..."
./scripts/test-new-feature.sh
```

### 4. Document in Scripts Table
Add entry to the scripts overview table above.

## Best Practices

### Script Development
- **Use Common Utilities**: Source `common.sh` for consistent behavior
- **Follow Naming**: Use `test-*.sh` pattern for test scripts
- **Error Handling**: Use `set -e` and proper exit codes
- **Logging**: Use `log_*()` functions for consistent output
- **Safety**: Default to dry-run mode, require explicit flags for real operations

### Performance Considerations
- **Timing**: Use `time_command()` for performance testing
- **Benchmarking**: Use `benchmark_command()` for repeated operations
- **Cleanup**: Always clean up test data and temporary files
- **Resource Limits**: Be mindful of memory and disk usage in tests

## CI/CD Integration

### GitHub Actions
```yaml
- name: Run Tests
  run: |
    ./scripts/start-test-env.sh
    ./scripts/integration-test.sh
```

### Local Development
```bash
# Watch mode (requires entr)
find . -name "*.go" | entr -c make test

# Full test suite
make test-all
```

### Makefile Targets
```bash
make test-integration    # Run integration tests
make test-performance   # Run performance tests
make test-all          # Run all test suites
```

## Troubleshooting

### Common Issues
- **Port conflicts**: Check `netstat -tulpn | grep :920[01]`
- **Container issues**: Run `./scripts/check-status.sh`
- **Config errors**: Verify `examples/test-config.yaml` exists
- **Permission errors**: Run `chmod +x scripts/*.sh`
- **Build failures**: Run `make clean && make build`

### Debug Commands
```bash
# Container logs
podman logs searchctl-elasticsearch
podman logs searchctl-opensearch

# Service health
curl localhost:9200/_cluster/health
curl localhost:9201/_cluster/health

# Config validation
./bin/searchctl config view

# Test specific functionality
./bin/searchctl --context elasticsearch --verbose cluster health
./bin/searchctl --context opensearch get indices --dry-run
```

### Performance Troubleshooting
```bash
# Check resource usage
podman stats

# Monitor test execution time
time ./scripts/integration-test.sh

# Run individual benchmarks
./scripts/test-performance.sh | grep "Results:"
```
