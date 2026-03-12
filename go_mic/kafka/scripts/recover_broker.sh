#!/usr/bin/env bash

BROKER=${1:-kafka1}
echo "=== Starting broker: $BROKER ==="
docker compose start "$BROKER"
echo "Broker $BROKER started at $(date)"
echo ""

# Ждём, пока брокер станет частью кластера
sleep 5
echo "Cluster status:"
docker compose exec kafka0 kafka-broker-api-versions \
  --bootstrap-server kafka0:29092 2>/dev/null | grep "id:" || \
  echo "(check manually with kafka-broker-api-versions)"
