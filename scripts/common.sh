#!/bin/bash

# Common utilities for searchctl test scripts
# Source this file at the beginning of test scripts for consistent functionality

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO] $1${NC}"
}

log_success() {
    echo -e "${GREEN}[SUCCESS] $1${NC}"
}

log_error() {
    echo -e "${RED}[ERROR] $1${NC}"
}

log_test() {
    echo -e "${YELLOW}[TEST] $1${NC}"
}

# Build searchctl and set up test environment
setup_test_environment() {
    export SEARCHCTL_CONFIG="examples/test-config.yaml"
    echo "[BUILD] Building searchctl..."
    make build
    log_success "Build completed"
}

# Run a command and capture timing information
time_command() {
    local cmd="$1"
    local description="$2"
    local start_time=$(date +%s.%N)
    
    echo "[TIMING] Starting: $description"
    eval "$cmd"
    local exit_code=$?
    
    local end_time=$(date +%s.%N)
    local duration=$(echo "$end_time - $start_time" | bc)
    
    if [ $exit_code -eq 0 ]; then
        echo "[TIMING] Completed in ${duration}s: $description"
    else
        echo "[TIMING] Failed after ${duration}s: $description"
    fi
    
    return $exit_code
}

# Run a test command with proper error handling
test_command() {
    local cmd="$1"
    local quiet="${2:-false}"
    
    if [[ "$quiet" != "true" ]]; then
        echo "[EXEC] Running: $cmd"
    fi
    
    if eval "$cmd"; then
        if [[ "$quiet" != "true" ]]; then
            echo "[EXEC] Command succeeded"
        fi
        return 0
    else
        echo "[EXEC] Command failed: $cmd"
        return 1
    fi
}

# Wait for a service to be ready
wait_for_service() {
    local service_name="$1"
    local url="$2"
    local max_attempts="${3:-30}"
    local attempt=1
    
    log_info "Waiting for $service_name to be ready..."
    
    while [ $attempt -le $max_attempts ]; do
        if curl -f "$url" >/dev/null 2>&1; then
            log_success "$service_name is ready"
            return 0
        fi
        
        echo "Attempt $attempt/$max_attempts: $service_name not ready yet..."
        sleep 2
        ((attempt++))
    done
    
    log_error "$service_name failed to start after $max_attempts attempts"
    return 1
}

# Clean up any test data or temporary files
cleanup_test_data() {
    log_info "Cleaning up test data..."
    # Add cleanup commands here as needed
    log_success "Cleanup completed"
}

# Check if test environment is running
check_environment() {
    log_info "Checking test environment..."
    
    log_info "Checking Elasticsearch on port 9200..."
    if ! curl -f http://localhost:9200/_cluster/health >/dev/null 2>&1; then
        log_error "Elasticsearch not running. Start test environment first:"
        log_error "   ./scripts/start-test-env.sh"
        exit 1
    fi
    log_success "Elasticsearch is ready"

    log_info "Checking OpenSearch on port 9201..."
    if ! curl -f http://localhost:9201/_cluster/health >/dev/null 2>&1; then
        log_error "OpenSearch not running. Start test environment first:"
        log_error "   ./scripts/start-test-env.sh"
        exit 1
    fi
    log_success "OpenSearch is ready"
    
    log_success "Test environment is ready"
}

# Test both engines with a command
test_both_engines() {
    local base_cmd="$1"
    local description="$2"
    
    log_test "Testing $description"
    
    # Test Elasticsearch
    log_info "Testing with Elasticsearch..."
    if test_command "${base_cmd/--context /--context elasticsearch }"; then
        log_success "Elasticsearch test passed"
    else
        log_error "Elasticsearch test failed"
        return 1
    fi
    
    # Test OpenSearch
    log_info "Testing with OpenSearch..."
    if test_command "${base_cmd/--context /--context opensearch }"; then
        log_success "OpenSearch test passed"
    else
        log_error "OpenSearch test failed"
        return 1
    fi
    
    log_success "$description completed successfully"
}

# Performance testing utilities
benchmark_command() {
    local cmd="$1"
    local iterations="${2:-10}"
    local description="$3"
    
    echo "[BENCHMARK] Starting: $description ($iterations iterations)"
    
    local total_time=0
    local successful_runs=0
    
    for i in $(seq 1 $iterations); do
        local start_time=$(date +%s.%N)
        
        if eval "$cmd" >/dev/null 2>&1; then
            local end_time=$(date +%s.%N)
            local duration=$(echo "$end_time - $start_time" | bc)
            total_time=$(echo "$total_time + $duration" | bc)
            ((successful_runs++))
            echo "  Run $i: ${duration}s [OK]"
        else
            echo "  Run $i: FAILED [FAIL]"
        fi
    done
    
    if [ $successful_runs -gt 0 ]; then
        local avg_time=$(echo "scale=3; $total_time / $successful_runs" | bc)
        echo "[BENCHMARK] Results: $successful_runs/$iterations successful, avg time: ${avg_time}s"
    else
        echo "[BENCHMARK] Results: No successful runs"
        return 1
    fi
}

# Validate JSON output
validate_json() {
    local json_output="$1"
    if echo "$json_output" | jq . >/dev/null 2>&1; then
        return 0
    else
        log_error "Invalid JSON output"
        return 1
    fi
}

# Validate YAML output
validate_yaml() {
    local yaml_output="$1"
    if echo "$yaml_output" | python3 -c "import yaml, sys; yaml.safe_load(sys.stdin)" >/dev/null 2>&1; then
        return 0
    else
        log_error "Invalid YAML output"
        return 1
    fi
}

# Print script usage/help
print_usage() {
    local script_name="$1"
    local description="$2"
    
    echo "Usage: $script_name"
    echo ""
    echo "Description: $description"
    echo ""
    echo "Prerequisites:"
    echo "  - Test environment must be running (./scripts/start-test-env.sh)"
    echo "  - searchctl must be built (make build)"
    echo ""
    echo "Environment Variables:"
    echo "  SEARCHCTL_CONFIG - Configuration file (default: examples/test-config.yaml)"
    echo ""
}

# Trap to ensure cleanup on exit
setup_cleanup_trap() {
    trap cleanup_test_data EXIT
}

# Create a standard index template for testing
create_test_index_template() {
    local context="$1"
    local template_name="$2"
    local index_pattern="$3"
    local port
    
    if [[ "$context" == "elasticsearch" ]]; then
        port=9200
    else
        port=9201
    fi
    
    log_info "Creating index template '$template_name' for $context..."
    
    curl -s -X PUT "localhost:$port/_index_template/$template_name" \
        -H "Content-Type: application/json" \
        -d '{
            "index_patterns": ["'$index_pattern'"],
            "data_stream": {},
            "template": {
                "settings": {
                    "number_of_shards": 1,
                    "number_of_replicas": 0
                },
                "mappings": {
                    "properties": {
                        "@timestamp": {"type": "date"},
                        "message": {"type": "text"}
                    }
                }
            }
        }' >/dev/null 2>&1 || true
}

# Delete a test index template
delete_test_index_template() {
    local context="$1"
    local template_name="$2"
    local port
    
    if [[ "$context" == "elasticsearch" ]]; then
        port=9200
    else
        port=9201
    fi
    
    curl -s -X DELETE "localhost:$port/_index_template/$template_name" >/dev/null 2>&1 || true
}

# Apply an index template from a YAML file using searchctl apply
apply_test_index_template() {
    local context="$1"
    local template_file="$2"
    local template_name="${3:-}"
    
    log_info "Applying index template from '$template_file' for $context..."
    
    if [[ ! -f "$template_file" ]]; then
        log_error "Template file '$template_file' not found"
        return 1
    fi
    
    # Apply the template using searchctl apply command
    if ./bin/searchctl --context "$context" apply -f "$template_file" >/dev/null 2>&1; then
        log_success "Index template applied successfully from $template_file"
        return 0
    else
        log_error "Failed to apply index template from $template_file"
        return 1
    fi
}

# Export functions for use in other scripts
export -f log_info log_success log_error log_test
export -f setup_test_environment time_command test_command
export -f wait_for_service cleanup_test_data check_environment
export -f test_both_engines benchmark_command
export -f validate_json validate_yaml print_usage setup_cleanup_trap
export -f create_test_index_template delete_test_index_template apply_test_index_template
