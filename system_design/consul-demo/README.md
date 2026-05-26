# Consul Service Discovery Demo

## Топология

```
CL (curl)
    │
    ├──► api-app:8081  ──[Service Query]──► Consul:8500 ──► web-app:8080
    │                                          │
    └──────────────────────────────────────────┘
                                        (Health Check каждые 10s)
```

## Запуск

```bash
docker compose up --build
```

## Демонстрация фич

### 1. Consul UI (Service Catalog)
Открой в браузере: http://localhost:8500

Увидишь: web и api сервисы с зелеными health check

---

### 2. Service Query через HTTP API
```bash
# Список всех сервисов в каталоге
curl http://localhost:8500/v1/catalog/services

# Найти healthy инстансы "web"
curl http://localhost:8500/v1/health/service/web?passing=true

# api-app делает service discovery и проксирует к web-app
curl http://localhost:8081/call
```

---

### 3. Health Check — убиваем web-app
```bash
# Останавливаем web-app
docker compose stop web-app

# Consul увидит падение через 10-20 секунд
# Проверяем: в UI статус web станет красным

# api-app теперь вернет 502 (no healthy instances)
curl http://localhost:8081/call

# Поднимаем обратно
docker compose start web-app
```

---

### 4. Подписка на изменения (Blocking Query)
```bash
# Consul поддерживает long-polling через ?index=N
# Получаем текущий индекс
INDEX=$(curl -s http://localhost:8500/v1/health/service/web | jq -r '.[0].Node.ModifyIndex // 0')

# Ждем любых изменений в сервисе web (блокирующий запрос)
curl "http://localhost:8500/v1/health/service/web?index=$INDEX&wait=30s"
```

---

### 5. DNS интерфейс
```bash
# Consul отвечает на DNS запросы для *.service.consul
# Из контейнера (dns: consul в compose):
docker compose exec api-app nslookup web.service.consul consul
```

## Порты
| Порт | Сервис |
|------|--------|
| 8500 | Consul UI + HTTP API |
| 8600/UDP | Consul DNS |
| 8080 | web-app |
| 8081 | api-app |
