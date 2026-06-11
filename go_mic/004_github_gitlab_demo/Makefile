.PHONY: lint fmt test build up down tag-module

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
	docker compose up --build -d

down:
	docker compose down -v

run: fmt lint test up

tag-module:
	@echo "git tag libs/logger/v0.1.0 && git push origin libs/logger/v0.1.0"

