FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY main.go .
RUN go mod init api-app && go build -o api-app .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/api-app .
EXPOSE 8081
CMD ["./api-app"]
