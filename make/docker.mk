up:
	docker compose up -d --build backend frontend db token-cleaner

down:
	docker compose down

deploy: migrate up

logs:
	docker compose logs -f backend frontend db token-cleaner

shell-front:
	docker compose exec frontend sh

shell-back:
	docker compose exec backend sh

shell-db:
	docker compose exec db sh

docker-clean:
	docker system prune

docker-clean-all:
	docker system prune -f -a --volumes
	docker builder prune -af

docker-clean-cache:
	docker builder prune

docker-clean-project: down
	@docker images -a --filter "label=com.docker.compose.project=setlist-pwa" -q | xargs -r docker rmi
	docker builder prune

.PHONY: up down logs shell-front shell-back shell-db docker-clean docker-clean-all docker-clean-cache docker-clean-project