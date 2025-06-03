dev-migrate:
	go run cmd/migrator/main.go up --development

dev-migrate-down:
	go run cmd/migrator/main.go down --development

dev-server:
	go run cmd/api/main.go --development