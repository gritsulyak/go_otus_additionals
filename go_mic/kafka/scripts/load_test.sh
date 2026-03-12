#!/usr/bin/env bash
# Непрерывно шлём платежи, считаем успехи и ошибки

SUCCESS=0
ERRORS=0
i=1

echo "Starting load test... Press Ctrl+C to stop"

while true; do
  RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" \
    -X POST http://localhost:8081/payments \
    -H "Content-Type: application/json" \
    -d "{\"id\":\"pay-$(date +%s%N)-${i}\",\"amount\":$(( RANDOM % 1000 + 1 ))}")

  if [ "$RESPONSE" -eq 201 ]; then
    SUCCESS=$((SUCCESS + 1))
  else
    ERRORS=$((ERRORS + 1))
    echo "[ERROR] HTTP $RESPONSE (total errors: $ERRORS)"
  fi

  echo -ne "Sent: $((SUCCESS + ERRORS)) | Success: $SUCCESS | Errors: $ERRORS\r"
  i=$((i + 1))
  sleep 0.1
done
