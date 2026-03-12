FROM golang:1.22-alpine AS builder

# gcc + musl-dev — обязательны для CGO на Alpine
# pkgconf — нужен для pkg-config при линковке librdkafka
RUN apk add --no-cache gcc musl-dev pkgconf

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# -tags musl — ключевой флаг для Alpine/musl libc
# -ldflags '-extldflags "-static"' — полностью статический бинарник
RUN CGO_ENABLED=1 go build \
    -tags musl \
    -ldflags '-extldflags "-static"' \
    -o service .

# Финальный образ — scratch, никаких рантаймов
FROM scratch
WORKDIR /app
COPY --from=builder /app/service .
COPY migrations/ ./migrations/
ENTRYPOINT ["./service"]
