#!/bin/bash
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

cleanup() {
    log_info "Cleaning up test resources..."
    
    # Delete test data streams
    for ds in test-simple test-logs test-metrics; do
        ./bin/searchctl delete datastream $ds -y 2>/dev/null || true
    done
    
    # Delete test templates
    for template in test-datastream-template logs-datastream-template metrics-datastream-template; do
        curl -s -X DELETE "localhost:9200/_index_template/$template" 2>/dev/null || true
    done
    
    log_success "Cleanup completed"
}

# Trap cleanup on exit
trap cleanup EXIT

log_info "Starting Data Stream Template Testing Workflow"

# Build searchctl
log_info "Building searchctl..."
make build >/dev/null 2>&1
log_success "Build completed"

# Test 1: Simple Data Stream Template
log_info "Test 1: Simple Data Stream Template"

log_info "Applying simple template..."
./bin/searchctl apply -f examples/simple-datastream-template.yaml

log_info "Creating simple data stream..."
./bin/searchctl create datastream test-simple

log_info "Verifying data stream creation..."
./bin/searchctl get datastreams test-simple

log_success "Simple data stream test completed"

# Test 2: Comprehensive Logs Template
log_info "Test 2: Comprehensive Logs Data Stream"

log_info "Applying logs template..."
./bin/searchctl apply -f examples/datastream-template.yaml

log_info "Creating logs data stream..."
./bin/searchctl create datastream test-logs

log_info "Checking data stream details..."
./bin/searchctl get datastreams test-logs -o yaml

# Test indexing a document
log_info "Testing document indexing..."
curl -s -X POST "localhost:9200/test-logs/_doc" \
    -H "Content-Type: application/json" \
    -d '{
        "@timestamp": "2024-01-01T12:00:00Z",
        "level": "INFO",
        "message": "Test log message",
        "service": {
            "name": "test-service",
            "version": "1.0.0"
        },
        "host": {
            "name": "test-host"
        }
    }' >/dev/null

log_success "Document indexed successfully"

# Test 3: Metrics Data Stream
log_info "Test 3: Metrics Data Stream"

log_info "Applying metrics template..."
./bin/searchctl apply -f examples/metrics-datastream-template.yaml

log_info "Creating metrics data stream..."
./bin/searchctl create datastream test-metrics

log_info "Testing metrics document..."
curl -s -X POST "localhost:9200/test-metrics/_doc" \
    -H "Content-Type: application/json" \
    -d '{
        "@timestamp": "2024-01-01T12:00:00Z",
        "metric": {
            "name": "cpu.usage",
            "type": "gauge",
            "unit": "percent",
            "value": 75.5
        },
        "host": {
            "name": "server-01"
        },
        "service": {
            "name": "system-monitor"
        }
    }' >/dev/null

log_success "Metrics test completed"

# Test 4: Rollover Operations
log_info "Test 4: Rollover Operations"

log_info "Testing manual rollover..."
./bin/searchctl rollover datastream test-logs --max-age 1s --max-docs 1

log_info "Testing rollover with conditions file..."
./bin/searchctl rollover datastream test-logs -f examples/rollover-conditions.json

log_success "Rollover tests completed"

# Test 5: Data Stream Management
log_info "Test 5: Data Stream Management"

log_info "Listing all data streams..."
./bin/searchctl get datastreams

log_info "Getting specific data stream details..."
./bin/searchctl get datastreams test-* -o table

log_info "Checking backing indices..."
./bin/searchctl get indices .ds-test-*

log_success "Management tests completed"

# Summary
log_success "All Data Stream Template Tests Completed Successfully!"

echo
log_info "Summary of created resources:"
echo "  • Data Streams: test-simple, test-logs, test-metrics"
echo "  • Templates: test-datastream-template, logs-datastream-template, metrics-datastream-template"
echo "  • Sample documents indexed in test streams"

echo
log_info "Try these commands to explore further:"
echo "  ./bin/searchctl get datastreams --help"
echo "  ./bin/searchctl rollover datastream --help"
echo "  ./bin/searchctl describe datastream test-logs"

echo
log_warning "Resources will be cleaned up automatically"