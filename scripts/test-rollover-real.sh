#!/bin/bash
set -e

# Source common utilities
source "$(dirname "$0")/common.sh"

echo "[TEST-REAL] Setting up real test environment for rollover operations..."

# Set up test environment
setup_test_environment

# Function to create index template for data streams
create_index_template() {
    local context=$1
    local template_name="test-logs-template"
    
    log_info "Creating index template for $context..."
    
    local template_body='{
        "index_patterns": ["test-logs-*"],
        "data_stream": {},
        "template": {
            "settings": {
                "number_of_shards": 1,
                "number_of_replicas": 0
            },
            "mappings": {
                "properties": {
                    "@timestamp": {
                        "type": "date"
                    },
                    "message": {
                        "type": "text"
                    },
                    "level": {
                        "type": "keyword"
                    }
                }
            }
        }
    }'
    
    # Create index template using curl directly
    local port
    if [[ "$context" == "elasticsearch" ]]; then
        port=9200
    else
        port=9201
    fi
    
    echo "Creating template via curl on port $port..."
    curl -X PUT "localhost:$port/_index_template/$template_name" \
        -H "Content-Type: application/json" \
        -d "$template_body" || echo "Template creation failed or already exists"
}

# Function to add test documents to data stream
add_test_documents() {
    local context=$1
    local datastream_name=$2
    
    log_info "Adding test documents to $datastream_name on $context..."
    
    local port
    if [[ "$context" == "elasticsearch" ]]; then
        port=9200
    else
        port=9201
    fi
    
    # Add some test documents
    for i in {1..10}; do
        local doc='{
            "@timestamp": "'$(date -u +%Y-%m-%dT%H:%M:%S.%3NZ)'",
            "message": "Test log message '$i'",
            "level": "info"
        }'
        
        curl -X POST "localhost:$port/$datastream_name/_doc" \
            -H "Content-Type: application/json" \
            -d "$doc" >/dev/null 2>&1 || echo "Failed to add document $i"
    done
    
    echo "Added 10 test documents to $datastream_name"
}

# Function to test real rollover operations
test_real_rollover() {
    local context=$1
    log_info "Testing real rollover operations for $context..."
    
    local test_datastream="test-logs-rollover-$context"
    
    # Create the data stream
    echo "Creating data stream: $test_datastream"
    ./bin/searchctl --context $context create datastream $test_datastream || echo "Data stream creation failed or already exists"
    
    # Add some test data
    add_test_documents $context $test_datastream
    
    # Wait a moment for documents to be indexed
    sleep 2
    
    # Show current data stream state
    echo "Current data stream state:"
    ./bin/searchctl --context $context get datastreams $test_datastream -o wide
    
    # Test rollover with very low document threshold to trigger rollover
    echo "Testing rollover with low document threshold..."
    ./bin/searchctl --context $context rollover datastream $test_datastream --max-docs 5 -o json
    
    # Show updated state
    echo "Data stream state after rollover:"
    ./bin/searchctl --context $context get datastreams $test_datastream -o wide
    
    # Test age-based rollover (this won't trigger immediately but shows the command works)
    echo "Testing age-based rollover..."
    ./bin/searchctl --context $context rollover datastream $test_datastream --max-age 1s -o yaml
    
    # Clean up - delete the test data stream
    echo "Cleaning up test data stream..."
    ./bin/searchctl --context $context delete datastream $test_datastream || echo "Cleanup failed"
}

# Check environment and run tests
check_environment

# Warning message
echo ""
log_error "WARNING: This script will create real data streams and indices!"
log_info "Make sure you understand the impact before proceeding."
read -p "Continue? (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Aborted."
    exit 1
fi

# Create index templates first
create_index_template "elasticsearch"
create_index_template "opensearch"

# Wait for templates to be available
sleep 2

# Test with real data
echo ""
log_info "Testing Elasticsearch with real rollover operations..."
test_real_rollover "elasticsearch"

echo ""
log_info "Testing OpenSearch with real rollover operations..."
test_real_rollover "opensearch"

echo ""
log_success "Real rollover testing completed!"
log_info "Note: Some operations may fail if the cluster doesn't support certain features."
log_info "Check the cluster logs if you encounter any issues."
