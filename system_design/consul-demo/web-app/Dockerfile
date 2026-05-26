FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY main.go .
RUN go mod init web-app && go build -o web-app .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/web-app .
EXPOSE 8080
CMD ["./web-app"]
