#!/bin/bash
set -e

echo "[STOP] Stopping SearchCtl Test Environment..."

# Stop and remove containers
podman-compose down -v

# Clean up any remaining containers
echo "[CLEANUP] Cleaning up remaining containers..."
podman container prune -f 2>/dev/null || true

echo "[SUCCESS] Test environment stopped and cleaned up"
