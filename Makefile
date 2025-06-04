build-api:
	go build -o bin/api/main cmd/api/main.go

build-migrator:
	go build -o bin/migrator/main cmd/migrator/main.go

dev-migrate:
	go run cmd/migrator/main.go up --development

dev-migrate-down:
	go run cmd/migrator/main.go down --development

dev-server:
	go run cmd/api/main.go --development


test-clearcache:
	go clean -testcache

test:
	go test ./internal/services

test-coverage:
	go test ./internal/services -cover

test-htmlcoverage:
	go test  ./internal/services -coverprofile=coverage.out
	go tool cover -html=coverage.out

mocks:
	mockery

db-up:
	docker-compose up -d

db-down:
	docker-compose down