MIGRATIONS_PATH = migrations
PSQL_DATABASE_NAME = book-shop
PSQL_URL ?= postgres://pguser:pgpass@localhost:5432/$(PSQL_DATABASE_NAME)?sslmode=disable

db_migrate:
	migrate -path=$(MIGRATIONS_PATH) -database=$(PSQL_URL) up
db_migrate_down:
	migrate -path=$(MIGRATIONS_PATH) -database=$(PSQL_URL) down
db_create:
	docker-compose exec postgres psql -U pguser --command='create database "$(PSQL_DATABASE_NAME)";'