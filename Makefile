include .env
export

up:
	docker-compose up -d --build

down:
	docker-compose down

logs:
	docker-compose logs -f

db-up:
	docker-compose up -d db

migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Erreur : L'argument 'name' est requis. Usage : make migrate-create name=mon_nom_de_migration"; \
		exit 1; \
	fi
	cd backend && migrate create -ext sql -dir db/migrations -seq $(name)

migrate-up:
	migrate -path ./backend/db/migrations -database "$(DB_URL_MIGRATE)" up

migrate-down:
	migrate -path ./backend/db/migrations -database "$(DB_URL_MIGRATE)" down

migrate-force:
	migrate -path ./backend/db/migrations -database "$(DB_URL_MIGRATE)" force $(version)

shell-front:
	docker-compose exec frontend sh

shell-back:
	docker-compose exec backend sh

shell-db:
	docker-compose exec db sh

db-connect:
	docker-compose exec db psql -U $(POSTGRES_USER) -d $(POSTGRES_DB)