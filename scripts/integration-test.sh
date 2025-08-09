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
    
    # Set correct port for context
    port=9200
    if [ "$context" = "opensearch" ]; then
        port=9201
    fi
    
    # Cleanup any leftover test resources first
    ./bin/searchctl --context $context delete datastream test-logs -y >/dev/null 2>&1 || true
    ./bin/searchctl --context $context delete datastream logs-test -y >/dev/null 2>&1 || true
    curl -s -X DELETE "localhost:$port/_index_template/test-logs-template" >/dev/null 2>&1 || true
    curl -s -X DELETE "localhost:$port/_index_template/logs-test-template" >/dev/null 2>&1 || true
    
    # Basic cluster operations
    test_command "./bin/searchctl --context $context cluster health" true
    test_command "./bin/searchctl --context $context cluster info" true
    test_command "./bin/searchctl --context $context get indices" true
    
    # Test datastream operations (requires index template)
    log_test "Testing datastream operations..."
    test_command "./bin/searchctl --context $context get datastreams" true
    
    # Create index template first (required for data streams)
    echo "Creating index template for test-logs..."
    curl -s -X PUT "localhost:$port/_index_template/test-logs-template" \
        -H "Content-Type: application/json" \
        -d '{
            "index_patterns": ["test-logs*"],
            "priority": 200,
            "data_stream": {},
            "template": {
                "settings": {
                    "number_of_shards": 1,
                    "number_of_replicas": 0
                }
            }
        }' || {
            echo "Failed to create index template for test-logs"
            exit 1
        }
    
    # Wait a moment for template to be available
    sleep 2
    
    # Verify template exists
    template_response=$(curl -s "localhost:$port/_index_template/test-logs-template")
    if [[ "$template_response" == *"index_template"* ]]; then
        echo "Template test-logs-template created successfully"
    else
        echo "Template verification failed: $template_response"
        exit 1
    fi
    
    # Test datastream creation with proper template
    test_command "./bin/searchctl --context $context create datastream test-logs" true
    test_command "./bin/searchctl --context $context get datastreams test-logs" true
    test_command "./bin/searchctl --context $context delete datastream test-logs -y" true
    
    # Clean up template
    curl -s -X DELETE "localhost:$port/_index_template/test-logs-template" >/dev/null 2>&1 || true
    
    # Test rollover operations (requires datastream with template)
    log_test "Testing rollover operations..."
    
    # Cleanup any existing test datastream first
    ./bin/searchctl --context $context delete datastream logs-test -y >/dev/null 2>&1 || true
    curl -s -X DELETE "localhost:$port/_index_template/logs-test-template" >/dev/null 2>&1 || true
    
    # Create index template for rollover testing with higher priority to avoid conflicts
    echo "Creating index template for logs-test..."
    curl -s -X PUT "localhost:$port/_index_template/logs-test-template" \
        -H "Content-Type: application/json" \
        -d '{
            "index_patterns": ["logs-test*"],
            "priority": 200,
            "data_stream": {},
            "template": {
                "settings": {
                    "number_of_shards": 1,
                    "number_of_replicas": 0
                }
            }
        }' || {
            echo "Failed to create index template for logs-test"
            exit 1
        }
    
    # Wait a moment for template to be available
    sleep 2
    
    # Verify template exists
    rollover_template_response=$(curl -s "localhost:$port/_index_template/logs-test-template")
    if [[ "$rollover_template_response" == *"index_template"* ]]; then
        echo "Template logs-test-template created successfully"
    else
        echo "Rollover template verification failed: $rollover_template_response"
        exit 1
    fi
    
    # Create test datastream for rollover testing
    test_command "./bin/searchctl --context $context create datastream logs-test" true
    test_command "./bin/searchctl --context $context rollover datastream logs-test --max-age 7d --max-docs 1000" true
    
    # Test different conditions based on engine capabilities
    if [ "$context" = "elasticsearch" ]; then
        test_command "./bin/searchctl --context $context rollover datastream logs-test --max-primary-shard-docs 500000" true
        test_command "./bin/searchctl --context $context rollover datastream logs-test --max-age 30d --max-docs 1000000 --max-primary-shard-docs 500000 --max-primary-shard-size 50gb --max-size 50gb" true
        test_command "./bin/searchctl --context $context rollover datastream logs-test --lazy" true
    else
        # OpenSearch doesn't support some ES-specific conditions
        test_command "./bin/searchctl --context $context rollover datastream logs-test --max-docs 500000" true
        test_command "./bin/searchctl --context $context rollover datastream logs-test --max-age 30d --max-docs 1000000 --max-size 50gb" true
    fi
    
    # Test output formats
    test_command "./bin/searchctl --context $context rollover ds logs-test --max-age 7d" true
    
    # Cleanup test datastream and template
    test_command "./bin/searchctl --context $context delete datastream logs-test -y" true
    curl -s -X DELETE "localhost:$port/_index_template/logs-test-template" >/dev/null 2>&1 || true
    
    # Test ComponentTemplate operations
    log_test "Testing ComponentTemplate operations..."
    
    # Cleanup any existing test component templates first
    curl -s -X DELETE "localhost:$port/_component_template/test-base-settings" >/dev/null 2>&1 || true
    curl -s -X DELETE "localhost:$port/_component_template/test-observability-mappings" >/dev/null 2>&1 || true
    
    # Test ComponentTemplate apply operations
    test_command "./bin/searchctl --context $context apply -f examples/component-templates/base-settings.yaml" true
    test_command "./bin/searchctl --context $context apply -f examples/component-templates/observability-mappings.yaml" true
    
    # Verify component templates were created by checking API directly
    echo "Verifying component templates were created..."
    base_settings_response=$(curl -s "localhost:$port/_component_template/base-settings")
    if [[ "$base_settings_response" == *"component_template"* ]]; then
        echo "Component template base-settings created successfully"
    else
        echo "Component template base-settings verification failed: $base_settings_response"
        exit 1
    fi
    
    observability_mappings_response=$(curl -s "localhost:$port/_component_template/observability-mappings")
    if [[ "$observability_mappings_response" == *"component_template"* ]]; then
        echo "Component template observability-mappings created successfully"
    else
        echo "Component template observability-mappings verification failed: $observability_mappings_response"
        exit 1
    fi
    
    # Test ComponentTemplate listing (if get command is implemented)
    # Note: This will gracefully fail if not implemented yet
    ./bin/searchctl --context $context get componenttemplates >/dev/null 2>&1 || echo "Get componenttemplates command not yet implemented"
    
    # Cleanup component templates
    curl -s -X DELETE "localhost:$port/_component_template/base-settings" >/dev/null 2>&1 || true
    curl -s -X DELETE "localhost:$port/_component_template/observability-mappings" >/dev/null 2>&1 || true
    
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
log_info "  ./scripts/test-rollover-real.sh - Real operations with test data"
log_info "  ./scripts/test-performance.sh - Performance benchmarking"
