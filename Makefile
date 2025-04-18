# ─────────── Settings ───────────
APP_NAME := url-shortener
GO_FILES := $(shell find . -name '*.go' -not -path "./vendor/*")
SWAG := github.com/swaggo/swag/cmd/swag

# ─────────── Commands ───────────

## Setup: Create local folders for bind mounts
setup:
	mkdir -p mongo-data redis-data

## Run locally without Docker
run:
	go run ./cmd/main.go

## Build Go binary
build:
	go build -o $(APP_NAME) ./cmd/main.go

## Run tests with coverage
test:
	go test ./... -coverprofile=coverage.out

## View coverage report in terminal
coverage:
	go tool cover -func=coverage.out

## View coverage report in browser
coverage-html:
	go tool cover -html=coverage.out

## Format code
fmt:
	go fmt ./...

## Tidy dependencies
tidy:
	go mod tidy

## Run linter (assumes golangci-lint is installed)
lint:
	golangci-lint run

## Generate Swagger docs
swagger:
	go install $(SWAG)@latest
	swag init --parseDependency --parseInternal

## Docker: Build and run containers
docker-up:
	docker-compose up --build

## Docker: Stop and remove containers
docker-down:
	docker-compose down

## Clean compiled files and coverage
clean:
	rm -f $(APP_NAME) coverage.out

.PHONY: run build test coverage coverage-html fmt tidy lint swagger docker-up docker-down clean setup
