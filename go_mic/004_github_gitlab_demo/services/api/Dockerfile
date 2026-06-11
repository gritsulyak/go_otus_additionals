FROM golang:1.25-alpine AS builder
WORKDIR /app

# Сначала go.mod для кэширования зависимостей
COPY libs/logger/go.mod ./libs/logger/go.mod
COPY services/api/go.mod services/api/go.mod
COPY services/api/go.sum services/api/go.sum

ARG CACHE_BUST=1
# Исходники logger — любое изменение инвалидирует кэш начиная отсюда
COPY libs/logger/ ./libs/logger/

WORKDIR /app/services/api
RUN go mod download

WORKDIR /app
COPY services/api/ ./services/api/

RUN cd services/api && \
    CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w -extldflags=-static" \
    -trimpath \
    -o /bin/api ./cmd

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /bin/api /api
EXPOSE 8080
ENTRYPOINT ["/api"]