include .env
export

up:
	docker-compose up -d --build backend frontend db

down:
	docker-compose down

logs:
	docker-compose logs -f backend frontend db

migrate:
	docker-compose up migrator

deploy: migrate up

db-up:
	docker-compose up -d db

shell-front:
	docker-compose exec frontend sh

shell-back:
	docker-compose exec backend sh

shell-db:
	docker-compose exec db sh

db-connect:
	docker-compose exec db psql -U $(POSTGRES_USER) -d $(POSTGRES_DB)

docker-clean:
	docker system prune

docker-clean-all:
	docker system prune -a

docker-clean-cache:
	docker builder prune

docker-clean-project: down
	@docker images -a --filter "label=com.docker.compose.project=setlist-pwa" -q | xargs -r docker rmi
	docker builder prune