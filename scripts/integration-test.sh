#!/bin/bash
set -e

# Load common utilities
source "$(dirname "$0")/common.sh"

log_test "Running SearchCtl Integration Tests..."

# Setup test environment and build
setup_test_environment

# Test basic functionality for both engines
for context in elasticsearch opensearch; do
    log_test "Testing $context..."
    
    # Basic cluster operations
    test_command "./bin/searchctl --context $context cluster health" true
    test_command "./bin/searchctl --context $context cluster info" true
    test_command "./bin/searchctl --context $context get indices" true
    
    # Test datastream operations (dry-run)
    log_test "Testing datastream operations..."
    test_command "./bin/searchctl --context $context get datastreams" true
    test_command "./bin/searchctl --context $context create datastream test-logs --dry-run" true
    test_command "./bin/searchctl --context $context delete datastream test-logs --dry-run" true
    
    # Test rollover operations (dry-run)
    log_test "Testing rollover operations..."
    test_command "./bin/searchctl --context $context rollover datastream logs-test --dry-run --max-age 7d --max-docs 1000" true
    test_command "./bin/searchctl --context $context rollover datastream logs-test --dry-run --max-primary-shard-docs 500000" true
    test_command "./bin/searchctl --context $context rollover datastream logs-test --dry-run --lazy --max-age 1d" true
    test_command "./bin/searchctl --context $context rollover datastream logs-test --dry-run --max-age 30d --max-docs 1000000 --max-primary-shard-docs 500000 --max-primary-shard-size 50gb --max-size 50gb" true
    
    # Test output formats
    test_command "./bin/searchctl --context $context rollover ds logs-test --dry-run --max-age 7d" true
    
    log_success "$context tests completed"
done

# Test help documentation
log_test "Testing help documentation..."
./bin/searchctl rollover --help >/dev/null
./bin/searchctl rollover datastream --help >/dev/null
./bin/searchctl create datastream --help >/dev/null
./bin/searchctl delete datastream --help >/dev/null

log_success "Integration tests completed successfully!"
log_info "Additional test scripts available:"
log_info "  ./scripts/test-rollover.sh - Comprehensive rollover testing"  
log_info "  ./scripts/test-rollover-real.sh - Real operations with test data"
log_info "  ./scripts/test-performance.sh - Performance benchmarking"
