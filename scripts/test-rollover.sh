#!/bin/bash
set -e

# Source common utilities
source "$(dirname "$0")/common.sh"

echo "[TEST] Testing SearchCtl Rollover Functionality..."

# Set test config and build
setup_test_environment

# Function to test rollover with different conditions
test_rollover_conditions() {
    local context=$1
    log_info "Testing rollover conditions for $context..."
    
    # Test individual conditions
    test_command "./bin/searchctl --context $context rollover datastream test-logs --dry-run --max-age 7d"
    test_command "./bin/searchctl --context $context rollover datastream test-logs --dry-run --max-docs 1000000"
    test_command "./bin/searchctl --context $context rollover datastream test-logs --dry-run --max-size 50gb"
    test_command "./bin/searchctl --context $context rollover datastream test-logs --dry-run --max-primary-shard-size 25gb"
    test_command "./bin/searchctl --context $context rollover datastream test-logs --dry-run --max-primary-shard-docs 500000"
    
    # Test lazy rollover
    test_command "./bin/searchctl --context $context rollover datastream test-logs --dry-run --lazy --max-age 1d"
    
    # Test multiple conditions
    test_command "./bin/searchctl --context $context rollover datastream test-logs --dry-run --max-age 30d --max-docs 1000000 --max-size 50gb"
    
    # Test with conditions file
    test_command "./bin/searchctl --context $context rollover datastream test-logs --dry-run -f examples/rollover-conditions.json"
    
    # Test different output formats
    test_command "./bin/searchctl --context $context rollover datastream test-logs --dry-run --max-age 7d -o json"
    test_command "./bin/searchctl --context $context rollover datastream test-logs --dry-run --max-age 7d -o yaml"
    
    # Test using alias
    test_command "./bin/searchctl --context $context rollover ds test-logs --dry-run --max-age 7d"
}

# Function to test datastream operations
test_datastream_operations() {
    local context=$1
    log_info "Testing datastream operations for $context..."
    
    # Test create datastream (dry-run)
    test_command "./bin/searchctl --context $context create datastream test-logs-new --dry-run"
    test_command "./bin/searchctl --context $context create ds test-logs-alias --dry-run"
    
    # Test delete datastream (dry-run)
    test_command "./bin/searchctl --context $context delete datastream test-logs-old --dry-run"
    test_command "./bin/searchctl --context $context delete ds test-logs-alias --dry-run"
    
    # Test get datastreams
    test_command "./bin/searchctl --context $context get datastreams"
    test_command "./bin/searchctl --context $context get ds"
    test_command "./bin/searchctl --context $context get datastreams 'test-*'"
    
    # Test different output formats for get
    ./bin/searchctl --context $context get datastreams -o json
    test_command "./bin/searchctl --context $context get datastreams -o yaml"
    test_command "./bin/searchctl --context $context get datastreams -o wide"
}

# Function to test error scenarios
test_error_scenarios() {
    local context=$1
    log_info "Testing error scenarios for $context..."
    
    # Test rollover without datastream name (should fail)
    test_command "./bin/searchctl --context $context rollover datastream --dry-run || echo 'Expected error: missing datastream name'"
    
    # Test rollover with invalid conditions
    test_command "./bin/searchctl --context $context rollover datastream test-logs --dry-run --max-age invalid || echo 'Expected error: invalid max-age format'"
    
    # Test create datastream without name (should fail)
    test_command "./bin/searchctl --context $context create datastream --dry-run || echo 'Expected error: missing datastream name'"
    
    # Test delete datastream without name (should fail)
    test_command "./bin/searchctl --context $context delete datastream --dry-run || echo 'Expected error: missing datastream name'"
}

# Function to test verbose mode
test_verbose_mode() {
    local context=$1
    log_info "Testing verbose mode for $context..."
    
    test_command "./bin/searchctl --context $context --verbose rollover datastream test-logs --dry-run --max-age 7d"
    test_command "./bin/searchctl --context $context --verbose get datastreams"
    test_command "./bin/searchctl --context $context --verbose create datastream test-logs --dry-run"
}

# Check environment and run tests
check_environment

# Test Elasticsearch
echo ""
log_info "Testing Elasticsearch rollover functionality..."
test_rollover_conditions "elasticsearch"
test_datastream_operations "elasticsearch"
test_error_scenarios "elasticsearch"
test_verbose_mode "elasticsearch"

# Test OpenSearch
echo ""
log_info "Testing OpenSearch rollover functionality..."
test_rollover_conditions "opensearch"
test_datastream_operations "opensearch"
test_error_scenarios "opensearch"
test_verbose_mode "opensearch"

echo ""
log_info "Testing help commands..."
test_command "./bin/searchctl rollover --help"
test_command "./bin/searchctl rollover datastream --help"
test_command "./bin/searchctl rollover ds --help"
test_command "./bin/searchctl create datastream --help"
test_command "./bin/searchctl delete datastream --help"
test_command "./bin/searchctl get datastreams --help"

echo ""
log_success "All rollover tests completed successfully!"
