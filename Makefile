.PHONY: docs up down logs build run semgrep lint prepare_test test
printf = @printf "%s\t\t%s\n"

.DEFAULT_GOAL := help

help:
	@echo -e "Commands available:\n"
	$(printf) "up" "start all services (api + web) via docker compose"
	$(printf) "down" "stop and remove all services"
	$(printf) "logs" "tail logs from all running services"
	$(printf) "run" "run the Go API locally (no docker)"
	$(printf) "build" "build the Go binary"
	$(printf) "semgrep" "run semgrep vulnerability scan"
	$(printf) "lint" "run golangci-lint"
	$(printf) "prepare_test" "generate mocks for tests"
	$(printf) "test" "run all tests"
	$(printf) "docs" "regenerate swagger docs"
	@echo -e "\n'help' is shown by default if you do not specify a command."

up:
	docker compose up --build -d

down:
	docker compose down

logs:
	docker compose logs -f

run:
	go run ./cmd/api/main.go

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/api ./cmd/api/main.go

semgrep:
	semgrep --config p/ci --config p/golang --error .

lint:
	golangci-lint run -v --timeout 5m

prepare_test:
	go generate tests/mocks_generator.go

test: prepare_test
	go test -v ./...

docs:
	swag init --dir cmd/api,internal/server,internal/dto