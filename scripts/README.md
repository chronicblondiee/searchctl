# Test Scripts

Automated testing for searchctl functionality with consistent logging and shared utilities.

## Quick Start

```bash
# Setup test environment
./scripts/start-test-env.sh

# Run tests
./scripts/integration-test.sh        # Basic functionality
./scripts/test-rollover.sh          # Rollover features (dry-run)
./scripts/test-rollover-real.sh     # Real operations (creates data)
./scripts/test-performance.sh       # Performance benchmarks

# Cleanup
./scripts/stop-test-env.sh
```

## Scripts Overview

| Script | Purpose | Safety | Duration | Usage |
|--------|---------|---------|----------|-------|
| `integration-test.sh` | Core functionality validation | ✅ Safe (dry-run) | ~30s | Daily CI |
| `test-rollover.sh` | Rollover feature testing | ✅ Safe (dry-run) | ~45s | Feature validation |
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
- `log_info()`, `log_success()`, `log_error()` - Consistent logging
- `setup_test_environment()` - Build and configure
- `check_environment()` - Verify services are running
- `test_command()` - Execute with error handling
- `time_command()` - Performance timing
- `benchmark_command()` - Load testing
- `test_both_engines()` - Test ES and OS together

### 3. Update Integration Script
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
    ./scripts/test-rollover.sh
```

### Local Development
```bash
# Watch mode (requires entr)
find . -name "*.go" | entr -c make test

# Full test suite
make test-all
```

## Troubleshooting

### Common Issues
- **Port conflicts**: Check `netstat -tulpn | grep :920[01]`
- **Container issues**: Run `./scripts/check-status.sh`
- **Config errors**: Verify `examples/test-config.yaml` exists
- **Permission errors**: Run `chmod +x scripts/*.sh`

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
```

## Troubleshooting

### Common Issues
- **Port conflicts**: Check `netstat -tulpn | grep :920[01]`
- **Container issues**: Run `./scripts/check-status.sh`
- **Config errors**: Verify `examples/test-config.yaml` exists
- **Permission errors**: Run `chmod +x scripts/*.sh`

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
```
