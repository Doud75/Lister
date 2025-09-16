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

test: test-up run-playwright test-down
	@echo "✅ E2E Tests finished. Report available in frontend/playwright-report/index.html"
	@# Sur macOS/Linux, cette commande peut ouvrir le rapport.
	@# open ./frontend/playwright-report/index.html || xdg-open ./frontend/playwright-report/index.html

test-up:
	@echo "--- Cleaning up previous test environment ---"
	@docker-compose -f docker-compose.test.yml --env-file .env.test down -v --remove-orphans
	@echo "--- Building and starting test environment (DB, Backend with seed, Frontend) ---"
	@docker-compose -f docker-compose.test.yml --env-file .env.test up --build -d
	@echo "--- Waiting for frontend to be healthy before running tests ---"
	@until curl -s -f http://localhost:4001 > /dev/null; do \
		echo "Waiting for frontend_test service on port 4001..."; \
		sleep 2; \
	done
	@echo "--- Test environment is ready ---"

run-playwright:
	@echo "--- Running Playwright tests ---"
	@cd frontend && PLAYWRIGHT_TEST_BASE_URL=http://localhost:4001 npx playwright test

test-down:
	@echo "--- Tearing down test environment ---"
	@docker-compose -f docker-compose.test.yml --env-file .env.test down -v --remove-orphans