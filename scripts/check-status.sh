#!/bin/bash

echo "🔍 Checking container status..."
podman ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"

echo ""
echo "🏥 Checking health status..."
for container in searchctl-elasticsearch searchctl-opensearch; do
    if podman ps --filter name=$container --format "{{.Names}}" | grep -q $container; then
        health=$(podman inspect $container --format "{{.State.Health.Status}}" 2>/dev/null || echo "no healthcheck")
        echo "$container: $health"
    else
        echo "$container: not running"
    fi
done

echo ""
echo "📋 Recent logs from OpenSearch:"
podman logs searchctl-opensearch --tail=10 2>/dev/null || echo "Container not found or not running"
