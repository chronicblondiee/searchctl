#!/bin/bash
set -e

# Test script for ComponentTemplate functionality
# This script tests the complete CRUD operations for ComponentTemplates

# Load common utilities
source "$(dirname "$0")/common.sh"

log_test "Testing ComponentTemplate CRUD operations..."

# Setup test environment and build
setup_test_environment

# Test with dry-run first
log_test "Testing ComponentTemplate with dry-run..."

# Test apply dry-run
test_command "./bin/searchctl apply --dry-run -f examples/component-templates/base-settings.yaml" true
test_command "./bin/searchctl apply --dry-run -f examples/component-templates/observability-mappings.yaml" true
test_command "./bin/searchctl apply --dry-run -f examples/component-templates/observability-mappings-opensearch.yaml" true

# Test get dry-run  
test_command "./bin/searchctl get component-templates --dry-run" true

# Test delete dry-run
test_command "./bin/searchctl delete component-template test-template --dry-run" true

log_success "Dry-run tests completed"

# Test basic functionality for both engines (if test environment is available)
if curl -f http://localhost:9200/_cluster/health >/dev/null 2>&1; then
    log_test "Testing ComponentTemplate operations with live cluster..."
    
    for context in elasticsearch; do
        log_test "Testing ComponentTemplate CRUD with $context..."
        
        # Set correct port for context
        port=9200
        if [ "$context" = "opensearch" ]; then
            port=9201
        fi
        
        # Cleanup any existing test component templates first
        curl -s -X DELETE "localhost:$port/_component_template/base-settings" >/dev/null 2>&1 || true
        curl -s -X DELETE "localhost:$port/_component_template/observability-mappings" >/dev/null 2>&1 || true
        
        # Test ComponentTemplate apply operations
        test_command "./bin/searchctl --context $context apply -f examples/component-templates/base-settings.yaml" true
        
        # Use different component templates based on the engine
        if [ "$context" = "elasticsearch" ]; then
            test_command "./bin/searchctl --context $context apply -f examples/component-templates/observability-mappings.yaml" true
        else
            test_command "./bin/searchctl --context $context apply -f examples/component-templates/observability-mappings-opensearch.yaml" true
        fi
        
        # Test ComponentTemplate get operations
        test_command "./bin/searchctl --context $context get component-templates" true
        test_command "./bin/searchctl --context $context get component-templates base-settings" true
        test_command "./bin/searchctl --context $context get component-templates observability-mappings" true
        
        # Test different output formats
        test_command "./bin/searchctl --context $context get component-templates -o json" true
        test_command "./bin/searchctl --context $context get component-templates -o yaml" true
        
        # Test ComponentTemplate delete operations
        test_command "./bin/searchctl --context $context delete component-template base-settings -y" true
        test_command "./bin/searchctl --context $context delete component-template observability-mappings -y" true
        
        log_success "$context ComponentTemplate CRUD tests completed"
    done
    
    log_success "Live cluster tests completed"
else
    log_info "No test cluster available, skipping live tests"
    log_info "To run live tests, start the test environment:"
    log_info "  ./scripts/start-test-env.sh"
fi

# Test help documentation
log_test "Testing help documentation..."
./bin/searchctl get component-templates --help >/dev/null
./bin/searchctl delete component-template --help >/dev/null

log_success "ComponentTemplate functionality tests completed successfully!"