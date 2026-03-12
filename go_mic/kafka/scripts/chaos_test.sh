#!/usr/bin/env bash
# Имитирует цикл: падение → пауза → восстановление → пауза
# Запустите вместе с load_test.sh и наблюдайте за логами

BROKERS=("kafka1" "kafka2")
FAIL_DURATION=20  # секунд брокер будет "упавшим"

echo "Starting chaos test. Load test should be running in another terminal."
echo "Press Ctrl+C to stop"

for BROKER in "${BROKERS[@]}"; do
  echo ""
  echo ">>> Failing broker: $BROKER"
  docker compose stop "$BROKER"
  echo "    $BROKER down. Waiting ${FAIL_DURATION}s..."
  sleep $FAIL_DURATION

  echo ">>> Recovering broker: $BROKER"
  docker compose start "$BROKER"
  echo "    $BROKER up. Waiting 15s for rebalance..."
  sleep 15
done

echo "Chaos test cycle complete."
