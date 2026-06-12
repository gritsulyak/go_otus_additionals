.PHONY: lint fmt test build up down security run tag-module

tidy:
	cd libs/logger && go mod tidy
	cd services/api && go mod tidy

fmt:
	gofmt -w libs/logger services/api

lint: tidy
	cd libs/logger && golangci-lint run ./...
	cd services/api && golangci-lint run ./...

test:
	cd libs/logger && go test -v -race ./...
	cd services/api && go test -v -race ./...

build:
	docker compose build

up:
	docker compose up -d --build

down:
	docker compose down -v

# govulncheck нужно установить один раз: go install golang.org/x/vuln/cmd/govulncheck@latest
security:
	cd libs/logger && govulncheck ./...
	cd services/api && govulncheck ./...

run: fmt lint security test up

# например: make tag-module TAG=v0.2.1
TAG ?= v0.2.0
tag-module:
	git tag libs/logger/$(TAG)
	git push origin libs/logger/$(TAG)