#!/bin/bash
set -e

echo "[SETUP] Starting SearchCtl Test Environment..."

# Check if podman-compose is available
if ! command -v podman-compose &> /dev/null; then
    echo "[ERROR] podman-compose not found. Please install it first."
    exit 1
fi

# Stop any existing containers
echo "[CLEANUP] Cleaning up existing containers..."
podman-compose down 2>/dev/null || true

# Start containers
echo "[DOCKER] Starting containers..."
podman-compose up -d

echo "[WAIT] Waiting for services to be healthy..."

# Wait for Elasticsearch with timeout
echo "Checking Elasticsearch..."
timeout=120
elapsed=0
while ! curl -f http://localhost:9200/_cluster/health >/dev/null 2>&1; do
    if [ $elapsed -ge $timeout ]; then
        echo "[ERROR] Elasticsearch failed to start within ${timeout}s"
        echo "[LOGS] Container logs:"
        podman logs searchctl-elasticsearch --tail=20
        exit 1
    fi
    printf '.'
    sleep 2
    elapsed=$((elapsed + 2))
done
echo "[SUCCESS] Elasticsearch ready"

# Wait for OpenSearch with timeout
echo "Checking OpenSearch..."
elapsed=0
while ! curl -f http://localhost:9201/_cluster/health >/dev/null 2>&1; do
    if [ $elapsed -ge $timeout ]; then
        echo "[ERROR] OpenSearch failed to start within ${timeout}s"
        echo "[LOGS] Container logs:"
        podman logs searchctl-opensearch --tail=20
        exit 1
    fi
    printf '.'
    sleep 2
    elapsed=$((elapsed + 2))
done
echo "[SUCCESS] OpenSearch ready"

echo "[SUCCESS] Test environment ready!"
echo ""
echo "Services available:"
echo "[URL] Elasticsearch:       http://localhost:9200"
echo "[URL] OpenSearch:          http://localhost:9201"
echo "[URL] Kibana:              http://localhost:5602"
echo "[URL] OpenSearch Dashboards: http://localhost:5601"
echo ""
echo "Test your searchctl commands:"
echo "  export SEARCHCTL_CONFIG=examples/test-config.yaml"
echo "  ./bin/searchctl --context elasticsearch cluster health"
echo "  ./bin/searchctl --context opensearch cluster health"
