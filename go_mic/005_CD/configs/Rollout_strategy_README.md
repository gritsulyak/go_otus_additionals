
# Пункт В: Демонстрация Стратегии развертывания (RollingUpdate) - c этим configmap

RollingUpdate хорошо виден именно с Prometheus, потому что у него есть `/api/v1/status/config` для проверки состояния на каждом поде.

## Полный сценарий демо RollingUpdate

### Подготовка — наблюдатель в отдельном терминале

```bash
# Терминал 1: смотрим как поды пересоздаются в реальном времени
watch -n 1 kubectl get pods -n demo-space
```

Начальное состояние (все 3 запущены):

```
NAME                      READY   STATUS    RESTARTS   AGE
my-app-7d9f4b6c8-4xkpz   2/2     Running   0          5m
my-app-7d9f4b6c8-n8vrt   2/2     Running   0          5m
my-app-7d9f4b6c8-qs7wl   2/2     Running   0          5m
```


***

### Шаг 1 — Смотрим текущую стратегию

```bash
kubectl describe deployment my-app -n demo-space | grep -A5 "Strategy"
```

```
StrategyType:           RollingUpdate
RollingUpdateStrategy:  0 max unavailable, 1 max surge
```


***

### Шаг 2 — Триггерим RollingUpdate (меняем версию образа)

```bash
kubectl set image deployment/my-app \
  web-app=prom/prometheus:v2.52.0 \
  -n demo-space
```


***

### Шаг 3 — Наблюдаем процесс в Терминале 1

Именно здесь видна работа `maxSurge: 1` + `maxUnavailable: 0`:

```
# Момент 1: создаётся +1 новый под (surge), старые все живы
# Момент 2: новый под прошёл readinessProbe — старый удаляется
my-app-7d9f4b6c8-qs7wl   0/2     Terminating   0          5m   ← убивается

# Момент 3: цикл повторяется для следующего пода
my-app-7d9f4b6c8-n8vrt   0/2     Terminating         0          5m

# Финал: все 3 поды на новой версии
NAME                      READY   STATUS    RESTARTS   AGE
my-app-9b3c7f1d2-xk9mw   2/2     Running   0          2m
my-app-9b3c7f1d2-r7pln   2/2     Running   0          1m
my-app-9b3c7f1d2-tt4sz   2/2     Running   0          30s
```

> **Ключевой момент**: `READY` никогда не падает ниже `3/3` — `maxUnavailable: 0` гарантирует нулевой даунтайм.

***

### Шаг 4 — Следим за rollout статусом (Терминал 2)

```bash
kubectl rollout status deployment/my-app -n demo-space
```

```
Waiting for deployment "my-app" rollout to finish: 1 out of 3 new replicas have been updated...
Waiting for deployment "my-app" rollout to finish: 2 out of 3 new replicas have been updated...
Waiting for deployment "my-app" rollout to finish: 1 old replicas are pending termination...
deployment "my-app" successfully rolled out
```


***

### Шаг 5 — История ревизий

```bash
kubectl rollout history deployment/my-app -n demo-space
```

```
REVISION  CHANGE-CAUSE
1         <none>    ← v2.51.0
2         <none>    ← v2.52.0
```


***

### Шаг 6 — Rollback

```bash
# Откат к предыдущей версии
kubectl rollout undo deployment/my-app -n demo-space

# Убеждаемся что вернулась старая версия образа
kubectl get pods -n demo-space -o jsonpath=\
'{range .items[*]}{.metadata.name}{"\t"}{.spec.containers[0].image}{"\n"}{end}'
```

```
my-app-7d9f4b6c8-nn2kp   prom/prometheus:v2.51.0
my-app-7d9f4b6c8-pw8xt   prom/prometheus:v2.51.0
my-app-7d9f4b6c8-zz3vl   prom/prometheus:v2.51.0
```


***

### Шаг 7 —  нулевой даунтайм

Запустить **до** `kubectl set image` и оставить работать во время rollout:

```bash
# Функция — всегда свежий running под
get_pod() {
  kubectl get pods -n demo-space -l app=my-app \
    --field-selector=status.phase=Running \
    -o jsonpath='{.items[0].metadata.name}' 2>/dev/null
}

# Цикл с динамическим обновлением POD каждую итерацию
while true; do
  POD=$(get_pod)

  if [ -z "$POD" ]; then
    echo "$(date +%T) → [NO PODS] FAIL"
  else
    RESULT=$(kubectl exec -n demo-space "$POD" -c web-app -- \
      wget -qO- http://localhost:9090/-/healthy 2>/dev/null \
      || echo "FAIL")
    echo "$(date +%T) [${POD}] → $RESULT"
  fi

  sleep 2
done
```

**Ожидаемый вывод во время rollout:**

```
10:05:01 → OK
10:05:02 → OK
10:05:03 → OK    ← в этот момент идёт пересоздание подов
10:05:04 → OK
10:05:05 → OK    ← ни одного FAIL
```

Нет ни одного `FAIL` — это и есть живое доказательство что `maxUnavailable: 0` + `readinessProbe` обеспечивают zero-downtime deployment.

