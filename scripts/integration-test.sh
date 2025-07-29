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
    run_with_context "$context" cluster health >/dev/null
    run_with_context "$context" cluster info >/dev/null
    run_with_context "$context" get indices >/dev/null
    
    # Test datastream operations (dry-run)
    log_test "Testing datastream operations..."
    run_with_context "$context" get datastreams >/dev/null
    run_with_context "$context" create datastream test-logs --dry-run >/dev/null
    run_with_context "$context" delete datastream test-logs --dry-run >/dev/null
    
    # Test rollover operations (dry-run)
    log_test "Testing rollover operations..."
    run_with_context "$context" rollover datastream logs-test --dry-run --max-age 7d --max-docs 1000 >/dev/null
    run_with_context "$context" rollover datastream logs-test --dry-run --max-primary-shard-docs 500000 >/dev/null
    run_with_context "$context" rollover datastream logs-test --dry-run --lazy --max-age 1d >/dev/null
    run_with_context "$context" rollover datastream logs-test --dry-run -f examples/rollover-conditions.json >/dev/null
    
    # Test output formats
    run_with_context "$context" rollover ds logs-test --dry-run --max-age 7d -o json >/dev/null
    run_with_context "$context" get datastreams -o wide >/dev/null
    
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
log_info "  📄 ./scripts/test-rollover.sh - Comprehensive rollover testing"  
log_info "  🧪 ./scripts/test-rollover-real.sh - Real operations with test data"
log_info "  🚀 ./scripts/test-performance.sh - Performance benchmarking"
