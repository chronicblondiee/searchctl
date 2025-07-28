#!/bin/bash
set -e

echo "🚀 Starting SearchCtl Test Environment..."

# Check if podman-compose is available
if ! command -v podman-compose &> /dev/null; then
    echo "❌ podman-compose not found. Please install it first."
    exit 1
fi

# Stop any existing containers
echo "🧹 Cleaning up existing containers..."
podman-compose down 2>/dev/null || true

# Start containers
echo "🐳 Starting containers..."
podman-compose up -d

echo "⏳ Waiting for services to be healthy..."

# Wait for Elasticsearch with timeout
echo "Checking Elasticsearch..."
timeout=120
elapsed=0
while ! curl -f http://localhost:9200/_cluster/health >/dev/null 2>&1; do
    if [ $elapsed -ge $timeout ]; then
        echo "❌ Elasticsearch failed to start within ${timeout}s"
        echo "📋 Container logs:"
        podman logs searchctl-elasticsearch --tail=20
        exit 1
    fi
    printf '.'
    sleep 2
    elapsed=$((elapsed + 2))
done
echo "✅ Elasticsearch ready"

# Wait for OpenSearch with timeout
echo "Checking OpenSearch..."
elapsed=0
while ! curl -f http://localhost:9201/_cluster/health >/dev/null 2>&1; do
    if [ $elapsed -ge $timeout ]; then
        echo "❌ OpenSearch failed to start within ${timeout}s"
        echo "📋 Container logs:"
        podman logs searchctl-opensearch --tail=20
        exit 1
    fi
    printf '.'
    sleep 2
    elapsed=$((elapsed + 2))
done
echo "✅ OpenSearch ready"

echo "🎉 Test environment ready!"
echo ""
echo "Services available:"
echo "📊 Elasticsearch:       http://localhost:9200"
echo "📊 OpenSearch:          http://localhost:9201"
echo "🖥️  Kibana:              http://localhost:5602"
echo "🖥️  OpenSearch Dashboards: http://localhost:5601"
echo ""
echo "Test your searchctl commands:"
echo "  export SEARCHCTL_CONFIG=examples/test-config.yaml"
echo "  ./bin/searchctl --context elasticsearch cluster health"
echo "  ./bin/searchctl --context opensearch cluster health"
