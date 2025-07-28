#!/bin/bash
set -e

echo "🛑 Stopping SearchCtl Test Environment..."

# Stop and remove containers
podman-compose down -v

# Clean up any remaining containers
echo "🧹 Cleaning up remaining containers..."
podman container prune -f 2>/dev/null || true

echo "✅ Test environment stopped and cleaned up"
