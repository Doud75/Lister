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

docker-clean:
	@echo "Nettoyage du système Docker (conteneurs arrêtés, réseaux inutilisés, images en suspens, caches de build)..."
	docker system prune

docker-clean-all:
	@echo "ATTENTION : Nettoyage AGRESSIF du système Docker..."
	@echo "Suppression de TOUS les conteneurs arrêtés, réseaux, images et caches de build inutilisés."
	docker system prune -a

docker-clean-cache:
	@echo "Nettoyage spécifique du cache de build Docker..."
	docker builder prune

docker-clean-project: down
	@echo "Suppression des anciennes images Docker créées par ce projet..."
	@docker images -a --filter "label=com.docker.compose.project=setlist-pwa" -q | xargs -r docker rmi
	@echo "Nettoyage du cache de build..."
	docker builder prune