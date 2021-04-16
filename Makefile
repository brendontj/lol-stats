 
BASE_DIR = $(shell pwd)
MIGRATE = docker run --rm -it -v "$(BASE_DIR)/pkg/db/sql/migration/migrations":"/migrations" --network host migrate/migrate -path=/migrations/ -database "postgres://postgres:postgres@localhost:5432/lolstats?sslmode=disable&timezone=UTC"
N_VERSION ?= $(N)

migrate:
	$(MIGRATE) up $(N_VERSION)

migrate-down:
	$(MIGRATE) down $(N_VERSION)

db-create:
	docker-compose exec \
		-e PGPASSWORD=postgres \
		db psql -h localhost -U postgres -c \
		"CREATE DATABASE lolstats"

db-close-all-connections:
	docker-compose exec \
		-e PGPASSWORD=postgres \
		db psql -h localhost -U postgres -c \
		"SELECT pg_terminate_backend(pg_stat_activity.pid) FROM pg_stat_activity WHERE pg_stat_activity.datname = '#{database}' AND pid <> pg_backend_pid();"

db-drop: db-close-all-connections
	docker-compose exec \
		-e PGPASSWORD=postgres \
		db psql -h localhost -U postgres -c \
		"DROP DATABASE lolstats"

db-reset: db-drop db-create migrate

.PHONY: migrate migrate-down db-create db-close-all-connections db-drop db-close-all-connections db-drop