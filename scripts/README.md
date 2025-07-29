# Test Scripts for SearchCtl Rollover and DataStream Features

This directory contains comprehensive test scripts for validating the rollover and datastream functionality in searchctl.

## Available Test Scripts

### 🔧 `integration-test.sh`
**Main integration test suite**
- Tests basic functionality across Elasticsearch and OpenSearch
- Validates core commands work correctly
- Includes basic rollover and datastream testing
- **Usage**: `./scripts/integration-test.sh`

### 🔄 `test-rollover.sh`
**Comprehensive rollover testing (dry-run mode)**
- Tests all rollover conditions and parameters
- Validates error handling and edge cases
- Tests output formats and verbose mode
- Safe to run (all operations in dry-run mode)
- **Usage**: `./scripts/test-rollover.sh`

### 🧪 `test-rollover-real.sh`
**Real rollover operations with test data**
- Creates actual index templates and data streams
- Performs real rollover operations with test data
- Tests with actual documents and indices
- **Warning**: Creates and deletes test data streams
- **Usage**: `./scripts/test-rollover-real.sh`

### 🚀 `test-performance.sh`
**Performance and stress testing**
- Measures command execution times
- Tests concurrent operations
- Stress tests with rapid command execution
- Provides performance benchmarks
- **Usage**: `./scripts/test-performance.sh`

### 📄 `test-conditions.sh`
**Rollover conditions file testing**
- Tests various condition file formats
- Validates JSON parsing and error handling
- Tests combinations of file and CLI parameters
- **Usage**: `./scripts/test-conditions.sh`

### 🏥 `start-test-env.sh`
**Test environment setup**
- Starts Elasticsearch and OpenSearch containers
- Waits for services to be healthy
- Required before running other tests
- **Usage**: `./scripts/start-test-env.sh`

### 🛑 `stop-test-env.sh`
**Test environment cleanup**
- Stops and removes test containers
- Cleans up resources
- **Usage**: `./scripts/stop-test-env.sh`

### ✅ `check-status.sh`
**Environment status check**
- Checks container health status
- Shows service availability
- **Usage**: `./scripts/check-status.sh`

## Test Features Covered

### Rollover Commands
- ✅ Basic rollover operations
- ✅ All condition types (age, docs, size, primary shard size/docs)
- ✅ Lazy rollover functionality
- ✅ Conditions file support
- ✅ Multiple output formats (table, json, yaml)
- ✅ Dry-run mode
- ✅ Verbose mode
- ✅ Alias support (`rollover ds`)

### DataStream Commands
- ✅ Create datastream operations
- ✅ Delete datastream operations
- ✅ List datastreams with patterns
- ✅ Multiple output formats
- ✅ Dry-run mode
- ✅ Alias support (`create ds`, `delete ds`, `get ds`)

### Error Handling
- ✅ Missing required arguments
- ✅ Invalid condition formats
- ✅ Non-existent files
- ✅ Invalid JSON in conditions files
- ✅ Network connectivity issues

### Performance Testing
- ✅ Command execution timing
- ✅ Concurrent operation support
- ✅ Stress testing with multiple rapid commands
- ✅ Memory and resource usage validation

## Prerequisites

### System Requirements

**Required Software:**
- **Go 1.21+** - For building searchctl binary
- **Container Runtime**: 
  - Docker with docker-compose, OR
  - Podman with podman-compose
- **Command Line Tools**:
  - `curl` - For API health checks and direct cluster communication
  - `jq` - For JSON parsing and validation in test scripts
  - `make` - For build automation (uses Makefile)

**System Resources:**
- **Memory**: Minimum 4GB RAM (2GB for containers + 2GB for host)
- **Storage**: 5GB free disk space for container images and test data
- **Network**: Ports 9200, 9201, 9300, 9301, 5601 must be available

### Installation Verification

```bash
# Check Go installation
go version  # Should show 1.21 or higher

# Check container runtime
docker --version && docker-compose --version
# OR
podman --version && podman-compose --version

# Check required tools
curl --version
jq --version
make --version

# Verify port availability
netstat -tuln | grep -E ':(9200|9201|9300|9301|5601)\s'
# Should show no existing bindings
```

### Container Runtime Setup

#### Docker Setup
```bash
# Install Docker (Ubuntu/Debian)
sudo apt update
sudo apt install docker.io docker-compose

# Add user to docker group (logout/login required)
sudo usermod -aG docker $USER

# Start Docker service
sudo systemctl enable docker
sudo systemctl start docker
```

#### Podman Setup (Alternative)
```bash
# Install Podman (Ubuntu/Debian)
sudo apt update
sudo apt install podman podman-compose

# Enable rootless containers
echo 'export DOCKER_HOST=unix:///run/user/$UID/podman/podman.sock' >> ~/.bashrc
```

### Initial Setup

1. **Clone and Build**:
   ```bash
   git clone https://github.com/chronicblondiee/searchctl.git
   cd searchctl
   
   # Build the binary
   make build
   # OR
   go build -o searchctl .
   ```

2. **Start Test Environment**:
   ```bash
   ./scripts/start-test-env.sh
   ```

3. **Verify Environment**:
   ```bash
   ./scripts/check-status.sh
   ```

### Troubleshooting Setup

#### Common Issues

**Port Conflicts:**
```bash
# Check what's using the ports
sudo netstat -tulpn | grep :9200
sudo netstat -tulpn | grep :9201

# Stop conflicting services
sudo systemctl stop elasticsearch  # If system ES is running
```

**Container Startup Issues:**
```bash
# Check container logs
docker logs searchctl-elasticsearch
docker logs searchctl-opensearch

# Check container status
docker ps -a

# Restart containers
./scripts/stop-test-env.sh
./scripts/start-test-env.sh
```

**Memory Issues:**
```bash
# Check available memory
free -h

# Reduce container memory (edit docker-compose.yml)
# Change: "ES_JAVA_OPTS=-Xms512m -Xmx512m"
# To:     "ES_JAVA_OPTS=-Xms256m -Xmx256m"
```

**Permission Issues:**
```bash
# Make scripts executable
chmod +x scripts/*.sh

# Fix Docker permissions (if needed)
sudo chown $USER:$USER /var/run/docker.sock
```

## Running Tests

### Quick Test (Recommended)
```bash
# Start test environment
./scripts/start-test-env.sh

# Run basic integration tests
./scripts/integration-test.sh

# Run comprehensive rollover tests (dry-run, safe)
./scripts/test-rollover.sh
```

### Comprehensive Testing
```bash
# Start test environment
./scripts/start-test-env.sh

# Run all test suites
./scripts/integration-test.sh
./scripts/test-rollover.sh
./scripts/test-conditions.sh
./scripts/test-performance.sh

# Optional: Test with real data (creates/deletes test data)
./scripts/test-rollover-real.sh
```

### Clean Up
```bash
# Stop test environment
./scripts/stop-test-env.sh

# Remove test containers and networks
docker-compose down --volumes --remove-orphans

# Clean up any test data files
rm -rf /tmp/searchctl-*

# Optional: Remove container images to free space
docker rmi docker.elastic.co/elasticsearch/elasticsearch:8.11.0
docker rmi opensearchproject/opensearch:2.11.0
docker rmi opensearchproject/opensearch-dashboards:2.11.0
```

### Advanced Configuration

#### CI/CD Integration
```bash
# Environment variables for automated testing
export SEARCHCTL_CONFIG="examples/test-config.yaml"
export SEARCHCTL_CONTEXT="elasticsearch"  # Force single context
export TEST_TIMEOUT="60s"  # Increase timeout for slow CI

# Run in headless mode
./scripts/start-test-env.sh --detach
./scripts/integration-test.sh
./scripts/test-rollover.sh
./scripts/stop-test-env.sh
```

#### Performance Testing Configuration
```bash
# For performance testing, allocate more resources
# Edit docker-compose.yml:
# elasticsearch:
#   environment:
#     - "ES_JAVA_OPTS=-Xms1g -Xmx1g"
# opensearch:
#   environment:
#     - "OPENSEARCH_JAVA_OPTS=-Xms1g -Xmx1g"

# Run performance tests
./scripts/test-performance.sh
```

#### Multiple Engine Testing
```bash
# Test against different engine versions
# Create custom docker-compose.override.yml:
services:
  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch:8.12.0
  opensearch:
    image: opensearchproject/opensearch:2.12.0

# Start with override
docker-compose -f docker-compose.yml -f docker-compose.override.yml up -d
```

## Test Configuration

### Configuration Files

#### `examples/test-config.yaml`
Primary configuration file used by all test scripts:

```yaml
apiVersion: v1
kind: Config
current-context: elasticsearch
contexts:
- name: elasticsearch
  context:
    cluster: elasticsearch-local
    user: default
- name: opensearch
  context:
    cluster: opensearch-local
    user: default
clusters:
- name: elasticsearch-local
  cluster:
    server: http://localhost:9200
    insecure-skip-tls-verify: true
- name: opensearch-local
  cluster:
    server: http://localhost:9201
    insecure-skip-tls-verify: true
users:
- name: default
  user: {}
```

#### `examples/rollover-conditions.json`
Sample rollover conditions file used in testing:

```json
{
  "max_age": "30d",
  "max_docs": 1000000,
  "max_size": "50gb",
  "max_primary_shard_size": "50gb",
  "max_primary_shard_docs": 500000
}
```

#### `docker-compose.yml`
Test environment definition with Elasticsearch and OpenSearch containers:

```yaml
services:
  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch:8.11.0
    ports: ["9200:9200", "9300:9300"]
    environment:
      - discovery.type=single-node
      - xpack.security.enabled=false
      - "ES_JAVA_OPTS=-Xms512m -Xmx512m"
  
  opensearch:
    image: opensearchproject/opensearch:2.11.0
    ports: ["9201:9200", "9301:9300"]
    environment:
      - discovery.type=single-node
      - "DISABLE_SECURITY_PLUGIN=true"
```

### Environment Variables

The test scripts use these environment variables:

- `SEARCHCTL_CONFIG`: Path to config file (default: `examples/test-config.yaml`)
- `SEARCHCTL_CONTEXT`: Override current context for single-context testing
- `TEST_TIMEOUT`: Maximum time to wait for operations (default: 30s)

### Test Environment Setup

#### 1. Container Configuration
```bash
# Start test environment (Elasticsearch + OpenSearch)
./scripts/start-test-env.sh

# Check services are healthy
./scripts/check-status.sh

# View container logs if needed
docker logs searchctl-elasticsearch
docker logs searchctl-opensearch
```

#### 2. Port Configuration
- **Elasticsearch**: `localhost:9200` (HTTP), `localhost:9300` (Transport)
- **OpenSearch**: `localhost:9201` (HTTP), `localhost:9301` (Transport)
- **OpenSearch Dashboards**: `localhost:5601` (Web UI)

#### 3. Security Configuration
Both services run with security disabled for testing:
- No authentication required
- TLS verification disabled
- Demo configurations disabled

### Custom Configuration

#### Creating Custom Test Config
```bash
# Copy and modify the test config
cp examples/test-config.yaml my-test-config.yaml

# Export to use custom config
export SEARCHCTL_CONFIG="my-test-config.yaml"

# Run tests with custom config
./scripts/integration-test.sh
```

#### Adding New Contexts
```yaml
# Add to your config file
contexts:
- name: my-cluster
  context:
    cluster: my-elasticsearch
    user: my-user
clusters:
- name: my-elasticsearch
  cluster:
    server: https://my-cluster.example.com:9200
    insecure-skip-tls-verify: false
users:
- name: my-user
  user:
    username: elastic
    password: changeme
```

#### Authentication Setup
```yaml
# Basic Auth
users:
- name: basic-user
  user:
    username: elastic
    password: changeme

# API Key Auth
users:
- name: api-user
  user:
    api-key: "base64-encoded-api-key"

# Certificate Auth
clusters:
- name: secure-cluster
  cluster:
    server: https://secure.example.com:9200
    certificate-authority: /path/to/ca.crt
    insecure-skip-tls-verify: false
```

## Expected Results

### Successful Test Run
- All commands execute without errors
- Dry-run operations show expected output
- Help commands display correct information
- Performance tests complete within reasonable time

### Common Issues
- **Connection Refused**: Test environment not started
- **Permission Denied**: Script not executable (`chmod +x scripts/*.sh`)
- **Template Errors**: Index templates not created (run `test-rollover-real.sh`)

## Continuous Integration

These scripts are designed to be CI-friendly:
- Return appropriate exit codes
- Provide clear success/failure indication
- Support automated execution
- Include timing and performance metrics

## Contributing

When adding new rollover or datastream features:
1. Add tests to the appropriate script
2. Update this README
3. Ensure all test modes are covered (dry-run, real operations, error cases)
4. Test with both Elasticsearch and OpenSearch
