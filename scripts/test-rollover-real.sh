#!/bin/bash
set -e

echo "🧪 Setting up real test environment for rollover operations..."

# Set test config
export SEARCHCTL_CONFIG="examples/test-config.yaml"

# Build searchctl
echo "🔨 Building searchctl..."
make build

# Function to create index template for data streams
create_index_template() {
    local context=$1
    local template_name="test-logs-template"
    
    echo "📝 Creating index template for $context..."
    
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
    
    echo "📄 Adding test documents to $datastream_name on $context..."
    
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
    echo "🔄 Testing real rollover operations for $context..."
    
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

# Check if test environment is running
echo "🏥 Checking test environment..."
if ! curl -f http://localhost:9200/_cluster/health >/dev/null 2>&1; then
    echo "❌ Elasticsearch not running. Start test environment first:"
    echo "   ./scripts/start-test-env.sh"
    exit 1
fi

if ! curl -f http://localhost:9201/_cluster/health >/dev/null 2>&1; then
    echo "❌ OpenSearch not running. Start test environment first:"
    echo "   ./scripts/start-test-env.sh"
    exit 1
fi

echo "✅ Test environment is ready"

# Create index templates first
create_index_template "elasticsearch"
create_index_template "opensearch"

# Wait for templates to be available
sleep 2

# Test with real data
echo ""
echo "🧪 Testing Elasticsearch with real rollover operations..."
test_real_rollover "elasticsearch"

echo ""
echo "🧪 Testing OpenSearch with real rollover operations..."
test_real_rollover "opensearch"

echo ""
echo "✅ Real rollover testing completed!"
echo "💡 Note: Some operations may fail if the cluster doesn't support certain features."
echo "📊 Check the cluster logs if you encounter any issues."
