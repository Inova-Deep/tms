#!/usr/bin/env bash
# =============================================================================
# TMS — Deploy update to VPS
# Run as root from /opt/tms: bash deploy/update.sh
# =============================================================================
set -euo pipefail

REPO_DIR="/opt/tms"
cd "$REPO_DIR"

echo "=== Pulling latest code ==="
git pull

echo "=== Rebuilding and restarting containers ==="
docker compose -f docker-compose.prod.yml up -d --build

echo "=== Cleaning up old images ==="
docker image prune -f

echo "=== Done! ==="
