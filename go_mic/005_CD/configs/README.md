# 0. Переменная с именем пода
POD=$(kubectl get pods -n demo-space -l app=my-app \
  -o jsonpath='{.items[0].metadata.name}')

# 1. Поды запущены
kubectl get pods -n demo-space

# 2. Конфиг ДО — запоминаем job_name
kubectl exec -n demo-space $POD -c web-app -- \
  wget -qO- http://localhost:9090/api/v1/status/config | grep job_name

# 3. Патчим ConfigMap — добавляем новый scrape job
kubectl patch configmap app-config -n demo-space --type merge -p \
  '{"data":{"prometheus.yml":"global:\n  scrape_interval: 15s\nscrape_configs:\n  - job_name: \"demo\"\n    static_configs:\n      - targets: [\"localhost:9090\"]\n  - job_name: \"demo-v2\"\n    static_configs:\n      - targets: [\"localhost:9091\"]\n"}}'

# 4. Ждём sync (смотрим лог релоадера)
kubectl logs -l app=my-app -c config-reloader -n demo-space -f

# 5. Конфиг ПОСЛЕ — видим demo-v2
kubectl exec -n demo-space $POD -c web-app -- \
  wget -qO- http://localhost:9090/api/v1/status/config | grep job_name

# 6. Финальное доказательство — метрика успеха
kubectl exec -n demo-space $POD -c web-app -- \
  wget -qO- 'http://localhost:9090/api/v1/query?query=prometheus_config_last_reload_successful'


