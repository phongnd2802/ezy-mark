GOOSE_DBSTRING ?= "user=root password=secret host=127.0.0.1 port=5400 dbname=daily-social sslmode=disable"
GOOSE_MIGRATION_DIR ?= internal/db/migrations
GOOSE_DRIVER ?= postgres

api:
	@go run cmd/api/main.go

server:
	@air


css:
	npx @tailwindcss/cli -i ./views/css/app.css -o ./public/css/style.css --watch

icons:
	@go run cmd/icongen/main.go


network:
	docker network create daily-social-network

postgres:
	docker run --name daily-social-db --network daily-social-network -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -p 5555:5432 postgres:15-alpine


createdb:
	docker exec -it postgres_container createdb --username=root --owner=root daily-social

dropdb:
	docker exec -it postgres_container dropdb daily-social


migrate-up:
	set GOOSE_DRIVER=$(GOOSE_DRIVER)
	set GOOSE_DBSTRING=$(GOOSE_DBSTRING)
	goose $(GOOSE_DRIVER) $(GOOSE_DBSTRING) -dir=$(GOOSE_MIGRATION_DIR) up

migrate-down:
	set GOOSE_DRIVER=$(GOOSE_DRIVER)
	set GOOSE_DBSTRING=$(GOOSE_DBSTRING)
	goose $(GOOSE_DRIVER) $(GOOSE_DBSTRING) -dir=$(GOOSE_MIGRATION_DIR) down-to 0

create-migration:
	set GOOSE_DRIVER=$(GOOSE_DRIVER)
	set GOOSE_DBSTRING=$(GOOSE_DBSTRING)
	goose -dir=$(GOOSE_MIGRATION_DIR) -s create $(name) sql


flush:
	docker exec -it redis_container redis-cli flushall

reset: flush migrate-down migrate-up sqlc

redis: 
	docker run --name redis-db --network daily-social-network  -p 6380:6379 -d redis:7-alpine

sqlc:
	sqlc generate

swag:
	swag init -g ./cmd/api/main.go -o ./docs

.PHONY: api flush server templ css icons network postgres migrate-up migrate-down create-migration sqlc createdb dropdb redis
.PHONY: swag