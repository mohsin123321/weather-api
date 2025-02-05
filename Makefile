.PHONY: docs
printf = @printf "%s\t\t%s\n"

.DEFAULT_GOAL := run

help:
	@echo -e "Commands available:\n"
	$(printf) "run" "execute the app"
	$(printf) "build" "build the app in an executable file"
	$(printf) "semgrep" "run semgrep to check for vulnerabilities"
	$(printf) "lint" "run the linter golangci-lint"
	$(printf) "prepare_test" "prepare mocks for the tests"
	$(printf) "test" "prepare tests and run the tests"
	$(printf) "docs" "generate the swagger documentation"

	@echo -e "\n'run' will be executed by default if you do not specify a command."

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