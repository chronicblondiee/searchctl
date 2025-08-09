#!/bin/bash
set -e

# Test script for LifecyclePolicy functionality
# This script tests the complete CRUD operations for LifecyclePolicies

# Load common utilities
source "$(dirname "$0")/common.sh"

log_test "Testing LifecyclePolicy CRUD operations..."

# Setup test environment and build
setup_test_environment

# Test with dry-run first
log_test "Testing LifecyclePolicy with dry-run..."

# Test apply dry-run
test_command "./bin/searchctl apply --dry-run -f examples/lifecycle-policies/basic-ilm-policy.yaml" true
test_command "./bin/searchctl apply --dry-run -f examples/lifecycle-policies/basic-ism-policy.yaml" true
test_command "./bin/searchctl apply --dry-run -f examples/lifecycle-policies/hot-warm-cold-policy.yaml" true
test_command "./bin/searchctl apply --dry-run -f examples/lifecycle-policies/delete-old-logs-policy.yaml" true

# Test get dry-run  
test_command "./bin/searchctl get lifecycle-policies --dry-run" true

# Test delete dry-run
test_command "./bin/searchctl delete lifecycle-policy test-policy --dry-run" true

log_success "Dry-run tests completed"

# Test basic functionality for both engines (if test environment is available)
if curl -f http://localhost:9200/_cluster/health >/dev/null 2>&1; then
    log_test "Testing LifecyclePolicy operations with live cluster..."
    
    for context in elasticsearch; do
        log_test "Testing LifecyclePolicy CRUD with $context..."
        
        # Set correct port for context
        port=9200
        if [ "$context" = "opensearch" ]; then
            port=9201
        fi
        
        # Cleanup any existing test lifecycle policies first
        curl -s -X DELETE "localhost:$port/_ilm/policy/basic-log-rotation" >/dev/null 2>&1 || true
        curl -s -X DELETE "localhost:$port/_plugins/_ism/policies/basic-log-rotation-ism" >/dev/null 2>&1 || true
        curl -s -X DELETE "localhost:$port/_ilm/policy/delete-old-logs" >/dev/null 2>&1 || true
        
        # Test LifecyclePolicy apply operations
        if [ "$context" = "elasticsearch" ]; then
            test_command "./bin/searchctl --context $context apply -f examples/lifecycle-policies/basic-ilm-policy.yaml" true
            test_command "./bin/searchctl --context $context apply -f examples/lifecycle-policies/delete-old-logs-policy.yaml" true
        else
            test_command "./bin/searchctl --context $context apply -f examples/lifecycle-policies/basic-ism-policy.yaml" true
        fi
        
        # Test LifecyclePolicy get operations
        test_command "./bin/searchctl --context $context get lifecycle-policies" true
        
        if [ "$context" = "elasticsearch" ]; then
            test_command "./bin/searchctl --context $context get lifecycle-policies basic-log-rotation" true
            test_command "./bin/searchctl --context $context get lifecycle-policies delete-old-logs" true
        else
            test_command "./bin/searchctl --context $context get lifecycle-policies basic-log-rotation-ism" true
        fi
        
        # Test different output formats
        test_command "./bin/searchctl --context $context get lifecycle-policies -o json" true
        test_command "./bin/searchctl --context $context get lifecycle-policies -o yaml" true
        
        # Test LifecyclePolicy delete operations
        if [ "$context" = "elasticsearch" ]; then
            test_command "./bin/searchctl --context $context delete lifecycle-policy basic-log-rotation -y" true
            test_command "./bin/searchctl --context $context delete lifecycle-policy delete-old-logs -y" true
        else
            test_command "./bin/searchctl --context $context delete lifecycle-policy basic-log-rotation-ism -y" true
        fi
        
        log_success "$context LifecyclePolicy CRUD tests completed"
    done
    
    log_success "Live cluster tests completed"
else
    log_info "No test cluster available, skipping live tests"
    log_info "To run live tests, start the test environment:"
    log_info "  ./scripts/start-test-env.sh"
fi

# Test help documentation
log_test "Testing help documentation..."
./bin/searchctl get lifecycle-policies --help >/dev/null
./bin/searchctl delete lifecycle-policy --help >/dev/null

# Test various aliases
log_test "Testing command aliases..."
./bin/searchctl get ilm --help >/dev/null
./bin/searchctl get ism --help >/dev/null
./bin/searchctl get lp --help >/dev/null
./bin/searchctl delete ilm test-policy --dry-run >/dev/null
./bin/searchctl delete ism test-policy --dry-run >/dev/null
./bin/searchctl delete lp test-policy --dry-run >/dev/null

log_success "LifecyclePolicy functionality tests completed successfully!"