# Test Scripts

Automated testing for searchctl functionality.

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

| Script | Purpose | Safety | Usage |
|--------|---------|---------|-------|
| `integration-test.sh` | Core functionality validation | ✅ Safe (dry-run) | Daily CI |
| `test-rollover.sh` | Rollover feature testing | ✅ Safe (dry-run) | Feature validation |
| `test-rollover-real.sh` | Real rollover operations | ⚠️ Creates data | Pre-release testing |
| `test-performance.sh` | Performance benchmarking | ✅ Safe (dry-run) | Performance regression |
| `test-conditions.sh` | Conditions file validation | ✅ Safe (dry-run) | Config testing |
| `start-test-env.sh` | Environment setup | ✅ Safe | Setup |
| `stop-test-env.sh` | Environment cleanup | ✅ Safe | Cleanup |
| `check-status.sh` | Health verification | ✅ Safe | Debugging |

## Test Categories

### Core Features
- **Cluster Operations**: health, info, connectivity
- **Resource Management**: indices, nodes, data streams
- **CRUD Operations**: create, get, delete resources

### Rollover Features  
- **Conditions**: age, docs, size, primary shard limits
- **Output Formats**: table, json, yaml
- **Advanced Options**: lazy rollover, conditions files
- **Error Handling**: missing args, invalid formats

### Data Stream Features
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
```bash
#!/bin/bash
set -e

# Common setup
source scripts/common.sh
setup_test_environment

# Your tests here
test_new_feature() {
    echo "Testing new feature..."
    ./bin/searchctl new-command --dry-run
}

# Run tests
test_new_feature
echo "✅ New feature tests passed"
```

### 2. Update Integration Script
Add new test calls to `integration-test.sh`:
```bash
echo "🧪 Testing new feature..."
./scripts/test-new-feature.sh
```

### 3. Document in README
Add entry to scripts table above.

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
