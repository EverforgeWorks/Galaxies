#!/bin/bash
# Move to project root (one level up)
cd "$(dirname "$0")/.."

echo "🚀 Stopping host services to prevent conflicts..."
sudo systemctl stop postgresql 2>/dev/null

echo "🐳 Building and Starting Galaxies..."
# We use the implicit docker-compose.yml in the root
sudo docker compose up --build -d

echo "✅ System Online."
sudo docker compose ps