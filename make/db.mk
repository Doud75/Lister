db-up:
	docker compose up -d db

migrate:
	docker compose up --build migrator

migrate-create:
	@if [ -z "$(name)" ]; then echo "Erreur: Il faut donner un nom"; exit 1; fi
	docker run --rm -v "$(shell pwd)/backend/db/migrations:/migrations" migrate/migrate create -ext sql -dir /migrations -seq $(name)

migrate-down:
	docker compose run --rm migrator /usr/local/bin/migrate -path=/migrations/ -database "postgres://postgres:password@lister-db-1:5432/setlist?sslmode=disable" down 1

db-connect:
	docker compose exec db psql -U $(POSTGRES_USER) -d $(POSTGRES_DB)

.PHONY: db-up migrate migrate-create db-connect