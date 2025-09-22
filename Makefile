-include .env
export

# --- General Development Commands ---
up:
	docker compose up -d --build backend frontend db

down:
	docker compose down

logs:
	docker compose logs -f backend frontend db

migrate:
	docker compose up migrator

deploy: migrate up


# --- Shell & DB Access ---
db-up:
	docker compose up -d db

shell-front:
	docker compose exec frontend sh

shell-back:
	docker compose exec backend sh

shell-db:
	docker compose exec db sh

db-connect:
	docker compose exec db psql -U $(POSTGRES_USER) -d $(POSTGRES_DB)


# --- Docker Cleanup Commands ---
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


# ==============================================================================
# --- TEST SUITE COMMANDS ---
# ==============================================================================

# --- Global Test Runner ---
# Lance tous les tests : d'abord les tests unitaires rapides, puis les tests E2E.
test-all: test-unit test
	@echo "✅ All tests (Unit & E2E) finished successfully."


# --- Unit Tests (Vitest) ---
# Rapides, sans Docker, pour les fonctions et composants isolés.
test-unit:
	@echo "--- Running Unit Tests ---"
	@cd frontend && npx vitest run
	@echo "✅ Unit Tests finished."

# Lance les tests unitaires en mode "watch" pour le développement.
test-unit-watch:
	@echo "--- Running Unit Tests in watch mode ---"
	@cd frontend && npx vitest


# --- End-to-End Tests (Playwright) ---
# Complets, avec Docker, pour les parcours utilisateurs.

# Lance TOUS les tests E2E
test: test-up run-playwright test-down docker-clean-all
	@echo "✅ E2E Tests finished. Report available in frontend/playwright-report/index.html"

# Tests E2E spécifiques par catégorie
test-setlist: test-up run-playwright-setlist test-down docker-clean-project
	@echo "✅ Setlist tests finished."

test-setlist-detail: test-up run-playwright-setlist-detail test-down docker-clean-project
	@echo "✅ Setlist detail tests finished."

test-setlist-add: test-up run-playwright-setlist-add test-down docker-clean-project
	@echo "✅ Setlist add tests finished."

test-setlist-new: test-up run-playwright-setlist-new test-down docker-clean-project
	@echo "✅ Setlist new tests finished."

test-song: test-up run-playwright-song test-down docker-clean-project
	@echo "✅ Song tests finished."

test-song-list: test-up run-playwright-song-list test-down docker-clean-project
	@echo "✅ Song list tests finished."

test-song-new: test-up run-playwright-song-new test-down docker-clean-project
	@echo "✅ Song new tests finished."

test-song-edit: test-up run-playwright-song-edit test-down docker-clean-project
	@echo "✅ Song edit tests finished."


# --- E2E Test Helpers (private commands) ---

test-up:
	@echo "--- Cleaning up previous test environment ---"
	@docker compose -f docker compose.test.yml --env-file .env.test down -v --remove-orphans
	@echo "--- Building and starting test environment (DB, Backend with seed, Frontend) ---"
	@docker compose -f docker compose.test.yml --env-file .env.test up --build -d
	@echo "--- Waiting for frontend to be healthy before running tests ---"
	@until curl -s -f http://localhost:4001 > /dev/null; do \
		echo "Waiting for frontend_test service on port 4001..."; \
		sleep 2; \
	done
	@echo "--- Test environment is ready ---"

run-playwright:
	@echo "--- Running ALL Playwright tests ---"
	@cd frontend && npx playwright test

run-playwright-setlist:
	@echo "--- Running SETLIST Playwright tests ---"
	@cd frontend && npx playwright test tests/setlist/

run-playwright-setlist-detail:
	@echo "--- Running SETLIST DETAIL Playwright tests ---"
	@cd frontend && npx playwright test tests/setlist/detail.spec.ts

run-playwright-setlist-add:
	@echo "--- Running SETLIST ADD Playwright tests ---"
	@cd frontend && npx playwright test tests/setlist/add.spec.ts

run-playwright-setlist-new:
	@echo "--- Running SETLIST NEW Playwright tests ---"
	@cd frontend && npx playwright test tests/setlist/new.spec.ts

run-playwright-song:
	@echo "--- Running SONG Playwright tests ---"
	@cd frontend && npx playwright test tests/song/

run-playwright-song-list:
	@echo "--- Running SONG LIST Playwright tests ---"
	@cd frontend && npx playwright test tests/song/list.spec.ts

run-playwright-song-new:
	@echo "--- Running SONG NEW Playwright tests ---"
	@cd frontend && npx playwright test tests/song/new.spec.ts

run-playwright-song-edit:
	@echo "--- Running SONG EDIT Playwright tests ---"
	@cd frontend && npx playwright test tests/song/edit.spec.ts

test-down:
	@echo "--- Tearing down test environment ---"
	@docker compose -f docker compose.test.yml --env-file .env.test down -v --remove-orphans