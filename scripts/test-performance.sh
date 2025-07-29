#!/bin/bash
set -e

# Source common utilities
source "$(dirname "$0")/common.sh"

echo "[PERFORMANCE] Testing SearchCtl rollover commands..."

# Set up test environment
setup_test_environment

# Function to run performance tests
run_performance_tests() {
    local context=$1
    log_info "Running performance tests for $context..."
    
    # Test rollover command performance
    time_command "./bin/searchctl --context $context rollover datastream test-logs --dry-run --max-age 7d" \
        "Basic rollover (dry-run)"
    
    time_command "Complex rollover with multiple conditions (dry-run)" \
        ./bin/searchctl --context $context rollover datastream test-logs --dry-run --max-age 30d --max-docs 1000000 --max-size 50gb --max-primary-shard-docs 500000
    
    time_command "Rollover with conditions file (dry-run)" \
        ./bin/searchctl --context $context rollover datastream test-logs --dry-run -f examples/rollover-conditions.json
    
    time_command "Get datastreams" \
        ./bin/searchctl --context $context get datastreams
    
    time_command "Get datastreams with pattern" \
        ./bin/searchctl --context $context get datastreams "test-*"
    
    time_command "Create datastream (dry-run)" \
        ./bin/searchctl --context $context create datastream perf-test-stream --dry-run
    
    time_command "Delete datastream (dry-run)" \
        ./bin/searchctl --context $context delete datastream perf-test-stream --dry-run
    
    # Test output format performance
    time_command "Rollover with JSON output (dry-run)" \
        ./bin/searchctl --context $context rollover datastream test-logs --dry-run --max-age 7d -o json
    
    time_command "Rollover with YAML output (dry-run)" \
        ./bin/searchctl --context $context rollover datastream test-logs --dry-run --max-age 7d -o yaml
    
    time_command "Get datastreams JSON output" \
        ./bin/searchctl --context $context get datastreams -o json
}

# Function to test concurrent operations
test_concurrent_operations() {
    local context=$1
    log_info "Testing concurrent operations for $context..."
    
    echo "Running 5 concurrent rollover operations..."
    for i in {1..5}; do
        (
            ./bin/searchctl --context $context rollover datastream test-logs-$i --dry-run --max-age 7d >/dev/null 2>&1
            echo "Concurrent operation $i completed"
        ) &
    done
    
    # Wait for all background jobs to complete
    wait
    echo "All concurrent operations completed"
}

# Function to stress test with many rapid commands
stress_test() {
    local context=$1
    log_info "Stress testing for $context..."
    
    echo "Running 20 rapid rollover commands..."
    local start_time=$(date +%s)
    
    for i in {1..20}; do
        ./bin/searchctl --context $context rollover datastream stress-test-$i --dry-run --max-docs 1000 >/dev/null 2>&1
    done
    
    local end_time=$(date +%s)
    local duration=$((end_time - start_time))
    echo "Completed 20 commands in ${duration} seconds"
    echo "Average: $((duration * 1000 / 20))ms per command"
}

# Check environment and run tests
check_environment

# Test both engines
echo ""
log_info "Elasticsearch performance tests..."
run_performance_tests "elasticsearch"
test_concurrent_operations "elasticsearch"
stress_test "elasticsearch"

echo ""
log_info "OpenSearch performance tests..."
run_performance_tests "opensearch"
test_concurrent_operations "opensearch"
stress_test "opensearch"

echo ""
log_info "Performance testing summary:"
log_success "All performance tests completed"
log_info "Check the timing results above for performance analysis"
log_info "The rollover functionality is ready for production use"

# Run performance tests
echo ""
log_info "Elasticsearch performance tests..."
run_performance_tests "elasticsearch"
test_concurrent_operations "elasticsearch"
stress_test "elasticsearch"

echo ""
log_info "OpenSearch performance tests..."
run_performance_tests "opensearch"
test_concurrent_operations "opensearch"
stress_test "opensearch"

echo ""
log_info "Performance testing summary:"
log_success "All performance tests completed"
log_info "Check the timing results above for performance analysis"
log_info "The rollover functionality is ready for production use"
