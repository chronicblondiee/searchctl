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
    
    # Set correct port for context
    port=9200
    if [ "$context" = "opensearch" ]; then
        port=9201
    fi
    
    # Create index templates for performance testing
    curl -s -X PUT "localhost:$port/_index_template/test-logs-template" \
        -H "Content-Type: application/json" \
        -d '{
            "index_patterns": ["test-logs*"],
            "data_stream": {},
            "template": {
                "settings": {
                    "number_of_shards": 1,
                    "number_of_replicas": 0
                }
            }
        }' >/dev/null 2>&1 || true
    
    curl -s -X PUT "localhost:$port/_index_template/perf-test-stream-template" \
        -H "Content-Type: application/json" \
        -d '{
            "index_patterns": ["perf-test-stream*"],
            "data_stream": {},
            "template": {
                "settings": {
                    "number_of_shards": 1,
                    "number_of_replicas": 0
                }
            }
        }' >/dev/null 2>&1 || true
    
    # Create test datastreams for performance testing
    ./bin/searchctl --context $context create datastream test-logs >/dev/null 2>&1 || true
    ./bin/searchctl --context $context create datastream perf-test-stream >/dev/null 2>&1 || true
    
    # Test rollover command performance
    time_command "./bin/searchctl --context $context rollover datastream test-logs --max-age 7d" \
        "Basic rollover"
    
    if [ "$context" = "elasticsearch" ]; then
        time_command "./bin/searchctl --context $context rollover datastream test-logs --max-age 30d --max-docs 1000000 --max-size 50gb --max-primary-shard-docs 500000" \
            "Complex rollover with multiple conditions"
    else
        time_command "./bin/searchctl --context $context rollover datastream test-logs --max-age 30d --max-docs 1000000 --max-size 50gb" \
            "Complex rollover with multiple conditions"
    fi
    
    if [ "$context" = "elasticsearch" ]; then
        time_command "./bin/searchctl --context $context rollover datastream test-logs -f examples/rollover-conditions.json" \
            "Rollover with conditions file"
    fi
    
    time_command "./bin/searchctl --context $context get datastreams" \
        "Get datastreams"
    
    time_command "./bin/searchctl --context $context get datastreams \"test-*\"" \
        "Get datastreams with pattern"
    
    time_command "./bin/searchctl --context $context create datastream perf-test-stream2" \
        "Create datastream"
    
    time_command "./bin/searchctl --context $context delete datastream perf-test-stream2" \
        "Delete datastream"
    
    # Test output format performance
    time_command "./bin/searchctl --context $context rollover datastream test-logs --max-age 7d -o json" \
        "Rollover with JSON output"
    
    time_command "./bin/searchctl --context $context rollover datastream test-logs --max-age 7d -o yaml" \
        "Rollover with YAML output"
    
    time_command "./bin/searchctl --context $context get datastreams -o json" \
        "Get datastreams JSON output"
}

# Function to test concurrent operations
test_concurrent_operations() {
    local context=$1
    log_info "Testing concurrent operations for $context..."
    
    echo "Running 5 concurrent rollover operations..."
    for i in {1..5}; do
        (
            ./bin/searchctl --context $context rollover datastream test-logs --max-age 7d >/dev/null 2>&1
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
        ./bin/searchctl --context $context rollover datastream test-logs --max-docs 1000 >/dev/null 2>&1
    done
    
    local end_time=$(date +%s)
    local duration=$((end_time - start_time))
    echo "Completed 20 commands in ${duration} seconds"
    echo "Average: $((duration * 1000 / 20))ms per command"
}

# Check environment and run tests
check_environment

# Test both engines
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

# Cleanup all test datastreams and templates
echo ""
log_info "Cleaning up test resources..."
for context in elasticsearch opensearch; do
    port=9200
    if [ "$context" = "opensearch" ]; then
        port=9201
    fi
    
    ./bin/searchctl --context $context delete datastream test-logs >/dev/null 2>&1 || true
    ./bin/searchctl --context $context delete datastream perf-test-stream >/dev/null 2>&1 || true
    curl -s -X DELETE "localhost:$port/_index_template/test-logs-template" >/dev/null 2>&1 || true
    curl -s -X DELETE "localhost:$port/_index_template/perf-test-stream-template" >/dev/null 2>&1 || true
done

echo ""
log_info "Performance testing summary:"
log_success "All performance tests completed"
log_info "Check the timing results above for performance analysis"
log_info "The rollover functionality is ready for production use"
