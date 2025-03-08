GOOSE_DBSTRING ?= "user=root password=secret host=127.0.0.1 port=5400 dbname=ezy-mark sslmode=disable"
GOOSE_MIGRATION_DIR ?= internal/db/migrations
GOOSE_DRIVER ?= postgres



dev:
	@go run cmd/server/main.go

docker-up:
	@docker-compose -f environment/docker-compose-dev.yml up

docker-down:
	@docker-compose -f environment/docker-compose-dev.yml down

network:
	docker network create ezy-mark-network

postgres:
	docker run --name ezy-mark-db --network ezy-mark-network -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -p 5555:5432 postgres:15-alpine

createdb:
	docker exec -it postgres_container createdb --username=root --owner=root ezy-mark

dropdb:
	docker exec -it postgres_container dropdb ezy-mark


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
	docker run --name redis-db --network ezy-mark-network  -p 6380:6379 -d redis:7-alpine

sqlc:
	sqlc generate

swag:
	swag init -g ./cmd/server/main.go -o ./docs

.PHONY: flush server network postgres migrate-up migrate-down create-migration sqlc createdb dropdb redis
.PHONY: swag docker-up docker-down