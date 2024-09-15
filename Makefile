codegen:
	mkdir -p internal/generated/schema
	oapi-codegen -package tender -generate chi-server,types,spec api/schema.yaml > internal/generated/schema/
up:
	docker-compose up -d
	go run cmd/main.go

run-migrate-development:
	sql-migrate up -env="development"

stop:
	docker compose stop