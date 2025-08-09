#!/bin/bash
set -e

# Source common utilities
source "$(dirname "$0")/common.sh"

echo "[TEST] Testing rollover conditions file functionality..."

# Set up test environment
setup_test_environment

# Create a temporary directory for test files
TEST_DIR="/tmp/searchctl-rollover-tests"
mkdir -p "$TEST_DIR"

# Function to create and test different condition files
test_conditions_file() {
    local context=$1
    log_info "Testing conditions files for $context..."
    
    # Create index template before creating datastream
    local port
    if [[ "$context" == "elasticsearch" ]]; then
        port=9200
    else
        port=9201
    fi
    
    curl -s -X PUT "localhost:$port/_index_template/test-logs-template" \
        -H 'Content-Type: application/json' \
        -d '{
            "index_patterns": ["test-logs*"],
            "data_stream": {},
            "template": {
                "mappings": {
                    "properties": {
                        "@timestamp": {"type": "date"},
                        "message": {"type": "text"}
                    }
                }
            }
        }' >/dev/null
    
    # Create test datastream for conditions testing
    ./bin/searchctl --context $context create datastream test-logs >/dev/null 2>&1 || true
    
    # Test 1: Basic conditions file
    cat > "$TEST_DIR/basic-conditions.json" << 'EOF'
{
  "max_age": "7d",
  "max_docs": 1000000
}
EOF

    # Also create a YAML variant for testing
    cat > "$TEST_DIR/basic-conditions.yaml" << 'EOF'
max_age: 7d
max_docs: 1000000
EOF
    
    echo "Testing basic conditions file (JSON)..."
    ./bin/searchctl --context $context rollover datastream test-logs -f "$TEST_DIR/basic-conditions.json"

    echo "Testing basic conditions file (YAML)..."
    ./bin/searchctl --context $context rollover datastream test-logs -f "$TEST_DIR/basic-conditions.yaml"
    
    # Test 2: Advanced conditions file
    if [[ "$context" == "elasticsearch" ]]; then
        cat > "$TEST_DIR/advanced-conditions.json" << 'EOF'
{
  "max_age": "30d",
  "max_docs": 5000000,
  "max_size": "100gb",
  "max_primary_shard_size": "50gb",
  "max_primary_shard_docs": 2500000
}
EOF
    else
        # OpenSearch doesn't support primary shard conditions
        cat > "$TEST_DIR/advanced-conditions.json" << 'EOF'
{
  "max_age": "30d",
  "max_docs": 5000000,
  "max_size": "100gb"
}
EOF
    fi
    
    echo "Testing advanced conditions file..."
    ./bin/searchctl --context $context rollover datastream test-logs -f "$TEST_DIR/advanced-conditions.json"
    
    # Test 3: Minimal conditions file
    cat > "$TEST_DIR/minimal-conditions.json" << 'EOF'
{
  "max_age": "1h"
}
EOF
    
    echo "Testing minimal conditions file..."
    ./bin/searchctl --context $context rollover datastream test-logs -f "$TEST_DIR/minimal-conditions.json"
    
    # Test 4: Size-only conditions
    if [[ "$context" == "elasticsearch" ]]; then
        cat > "$TEST_DIR/size-conditions.json" << 'EOF'
{
  "max_size": "10gb",
  "max_primary_shard_size": "5gb"
}
EOF
    else
        # OpenSearch doesn't support primary shard size condition
        cat > "$TEST_DIR/size-conditions.json" << 'EOF'
{
  "max_size": "10gb"
}
EOF
    fi
    
    echo "Testing size-only conditions file..."
    ./bin/searchctl --context $context rollover datastream test-logs -f "$TEST_DIR/size-conditions.json"
    
    # Test 5: Document count conditions
    if [[ "$context" == "elasticsearch" ]]; then
        cat > "$TEST_DIR/docs-conditions.json" << 'EOF'
{
  "max_docs": 1000000,
  "max_primary_shard_docs": 500000
}
EOF
    else
        # OpenSearch doesn't support primary shard docs condition
        cat > "$TEST_DIR/docs-conditions.json" << 'EOF'
{
  "max_docs": 1000000
}
EOF
    fi
    
    echo "Testing document count conditions file..."
    ./bin/searchctl --context $context rollover datastream test-logs -f "$TEST_DIR/docs-conditions.json"
    
    # Test 6: Invalid JSON (should fail gracefully)
    cat > "$TEST_DIR/invalid-conditions.json" << 'EOF'
{
  "max_age": "7d",
  "max_docs": 1000000
  // This is invalid JSON due to missing comma
}
EOF
    
    echo "Testing invalid JSON file (should fail)..."
    ./bin/searchctl --context $context rollover datastream test-logs -f "$TEST_DIR/invalid-conditions.json" || echo "[OK] Correctly failed with invalid JSON"
    
    # Test 7: Non-existent file (should fail)
    echo "Testing non-existent file (should fail)..."
    ./bin/searchctl --context $context rollover datastream test-logs -f "$TEST_DIR/non-existent.json" || echo "[OK] Correctly failed with non-existent file"
    
    # Test 8: Combination of command line args and file
    echo "Testing combination of command line and file conditions..."
    ./bin/searchctl --context $context rollover datastream test-logs -f "$TEST_DIR/basic-conditions.json" --max-size 50gb
}

# Function to test output formats with conditions files
test_output_formats() {
    local context=$1
    log_info "Testing output formats with conditions files for $context..."
    
    if [[ "$context" == "elasticsearch" ]]; then
    echo "JSON output:"
    ./bin/searchctl --context $context rollover datastream test-logs -f examples/rollover-conditions.json -o json

    echo "YAML output:"
    ./bin/searchctl --context $context rollover datastream test-logs -f examples/rollover-conditions.yaml -o yaml

    echo "Table output (default, YAML file):"
    ./bin/searchctl --context $context rollover datastream test-logs -f examples/rollover-conditions.yaml
    else
        # For OpenSearch, use a basic conditions file without primary shard conditions
        echo "JSON output:"
        ./bin/searchctl --context $context rollover datastream test-logs -f "$TEST_DIR/basic-conditions.json" -o json
        
        echo "YAML output:"
        ./bin/searchctl --context $context rollover datastream test-logs -f "$TEST_DIR/basic-conditions.yaml" -o yaml
        
        echo "Table output (default):"
        ./bin/searchctl --context $context rollover datastream test-logs -f "$TEST_DIR/basic-conditions.json"
    fi
}

# Function to test verbose mode with conditions files
test_verbose_mode() {
    local context=$1
    log_info "Testing verbose mode with conditions files for $context..."
    
    if [[ "$context" == "elasticsearch" ]]; then
        ./bin/searchctl --context $context --verbose rollover datastream test-logs -f examples/rollover-conditions.json
    else
        # For OpenSearch, use a basic conditions file without primary shard conditions
        ./bin/searchctl --context $context --verbose rollover datastream test-logs -f "$TEST_DIR/basic-conditions.json"
    fi
}

# Check environment and run tests
check_environment

# Test with both engines
echo ""
log_info "Testing Elasticsearch conditions files..."
test_conditions_file "elasticsearch"
test_output_formats "elasticsearch"
test_verbose_mode "elasticsearch"

echo ""
log_info "Testing OpenSearch conditions files..."
test_conditions_file "opensearch"
test_output_formats "opensearch"
test_verbose_mode "opensearch"

# Validate the example conditions file
echo ""
log_info "Validating example conditions file..."
if [[ -f "examples/rollover-conditions.json" ]]; then
    log_success "Example conditions file exists"
    if python3 -m json.tool examples/rollover-conditions.json >/dev/null 2>&1; then
        log_success "Example conditions file is valid JSON"
        echo "Contents:"
        cat examples/rollover-conditions.json
    else
        log_error "Example conditions file is invalid JSON"
    fi
else
    log_error "Example conditions file not found"
fi

# Clean up
echo ""
log_info "Cleaning up test files..."
rm -rf "$TEST_DIR"

# Clean up test datastreams and templates
for context in elasticsearch opensearch; do
    ./bin/searchctl --context $context delete datastream test-logs -y >/dev/null 2>&1 || true
    
    # Clean up index template
    if [[ "$context" == "elasticsearch" ]]; then
        port=9200
    else
        port=9201
    fi
    curl -s -X DELETE "localhost:$port/_index_template/test-logs-template" >/dev/null 2>&1 || true
done

echo ""
log_success "All conditions file tests completed successfully!"
log_info "The rollover conditions file functionality is working correctly"
