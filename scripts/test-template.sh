#!/bin/bash
# Template for new test scripts
set -e

# Load common utilities
source "$(dirname "$0")/common.sh"

# Script description
SCRIPT_NAME="New Feature Test"
SCRIPT_PURPOSE="Template for testing new features"

log_test "Running $SCRIPT_NAME..."
log_info "$SCRIPT_PURPOSE"

# Setup test environment
setup_test_environment

# Define test functions for your new feature
test_new_feature_basic() {
    local context="$1"
    log_test "Testing basic new feature functionality for $context..."
    
    # Add your basic feature tests here
    # Example:
    # run_with_context "$context" new-command --dry-run
    # run_with_context "$context" new-command list
    
    log_success "Basic new feature tests passed for $context"
}

test_new_feature_advanced() {
    local context="$1"
    log_test "Testing advanced new feature functionality for $context..."
    
    # Add your advanced feature tests here
    # Example:
    # run_with_context "$context" new-command --advanced-option
    # run_with_context "$context" new-command -o json
    
    log_success "Advanced new feature tests passed for $context"
}

test_new_feature_error_handling() {
    local context="$1"
    log_test "Testing new feature error handling for $context..."
    
    # Test error scenarios
    # Example:
    # test_command "New command without required args should fail" 1 \
    #     run_with_context "$context" new-command
    
    log_success "Error handling tests passed for $context"
}

# Run tests for both engines
for context in elasticsearch opensearch; do
    log_test "Testing new feature with $context..."
    
    test_new_feature_basic "$context"
    test_new_feature_advanced "$context"
    test_new_feature_error_handling "$context"
    
    log_success "All $context tests passed"
done

# Performance testing (optional)
if [ "${RUN_PERFORMANCE_TESTS:-false}" = "true" ]; then
    log_performance "Running performance tests..."
    
    for context in elasticsearch opensearch; do
        time_command "New feature performance test for $context" \
            run_with_context "$context" new-command --performance-test
    done
fi

log_success "$SCRIPT_NAME completed successfully!"
log_info "All new feature functionality is working correctly"
