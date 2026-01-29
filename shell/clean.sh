#!/bin/bash

echo "🧹 Starting Docker Cleanup..."

# 1. Remove dangling images (layers from old builds)
# -f forces it without prompt
docker image prune -f

# 2. OPTIONAL: Remove stopped containers and unused networks
# This is safe to run while the game is up.
docker container prune -f
docker network prune -f

echo "✨ Disk space reclaimed:"
docker system df
