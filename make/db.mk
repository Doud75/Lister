db-up:
	docker compose up -d db

migrate:
	docker compose up --build migrator

migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Erreur : L'argument 'name' est requis. Usage : make migrate-create name=mon_nom_de_migration"; \
		exit 1; \
	fi
	cd backend && migrate create -ext sql -dir db/migrations -seq $(name)

db-connect:
	docker compose exec db psql -U $(POSTGRES_USER) -d $(POSTGRES_DB)

.PHONY: db-up migrate migrate-create db-connect