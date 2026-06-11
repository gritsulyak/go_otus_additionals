FROM golang:1.25-alpine AS builder
WORKDIR /app

# Копируем go.mod/go.sum обоих модулей сохраняя структуру путей
COPY libs/logger/go.mod ./libs/logger/go.mod
COPY services/api/go.mod services/api/go.mod
COPY services/api/go.sum services/api/go.sum

WORKDIR /app/services/api
RUN go mod download

# Копируем исходники
WORKDIR /app
COPY libs/logger/ ./libs/logger/
COPY services/api/ ./services/api/

RUN cd services/api && CGO_ENABLED=0 go build -ldflags="-s -w" -o /bin/api ./cmd

FROM alpine:latest
RUN adduser -D -g '' appuser
USER appuser
COPY --from=builder /bin/api /api
EXPOSE 8080
ENTRYPOINT ["/api"]