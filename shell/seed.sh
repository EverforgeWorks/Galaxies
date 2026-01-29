#!/bin/bash
cd "$(dirname "$0")/.."

echo "🌱 Seeding Universe..."
# Run the compiled binary
sudo docker compose exec backend ./seeder
