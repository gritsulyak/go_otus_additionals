#!/usr/bin/env bash
# Останавливаем один брокер для имитации отказа

BROKER=${1:-kafka1}
echo "=== Stopping broker: $BROKER ==="
docker compose stop "$BROKER"
echo "Broker $BROKER stopped at $(date)"
echo ""
echo "Remaining brokers:"
docker compose ps kafka0 kafka1 kafka2 | grep Up

echo ""
echo "Watch logs of payment-service for reconnect events:"
echo "  docker compose logs -f payment-service"
echo ""
echo "To recover, run: ./scripts/recover_broker.sh $BROKER"
