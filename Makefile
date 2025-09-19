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
	docker system prune -f -a --volumes

docker-clean-cache:
	docker builder prune

docker-clean-project: down
	@docker images -a --filter "label=com.docker.compose.project=setlist-pwa" -q | xargs -r docker rmi
	docker builder prune



# --- Test Commands ---

# Lance TOUS les tests E2E
test: test-up run-playwright test-down docker-clean-all
	@echo "✅ E2E Tests finished. Report available in frontend/playwright-report/index.html"
	@# Sur macOS/Linux, cette commande peut ouvrir le rapport.
	@# open ./frontend/playwright-report/index.html || xdg-open ./frontend/playwright-report/index.html

# Lance tous les tests relatifs aux SETLISTS
test-setlist: test-up run-playwright-setlist test-down docker-clean-all
	@echo "✅ Setlist tests finished. Report available in frontend/playwright-report/index.html"

# Lance uniquement les tests de la page de DÉTAIL d'une setlist
test-setlist-detail: test-up run-playwright-setlist-detail test-down docker-clean-all
	@echo "✅ Setlist detail tests finished. Report available in frontend/playwright-report/index.html"

# Lance uniquement les tests de la page d'AJOUT d'une setlist
test-setlist-add: test-up run-playwright-setlist-add test-down docker-clean-all
	@echo "✅ Setlist add tests finished. Report available in frontend/playwright-report/index.html"

# Lance uniquement les tests de la page de CRÉATION d'une setlist
test-setlist-new: test-up run-playwright-setlist-new test-down docker-clean-all
	@echo "✅ Setlist new tests finished. Report available in frontend/playwright-report/index.html"

test-song: test-up run-playwright-song test-down docker-clean-all
	@echo "✅ Song tests finished. Report available in frontend/playwright-report/index.html"

# Lance uniquement les tests de la page LISTE des chansons
test-song-list: test-up run-playwright-song-list test-down docker-clean-all
	@echo "✅ Song list tests finished. Report available in frontend/playwright-report/index.html"

# Lance uniquement les tests de la page de CRÉATION d'une chanson
test-song-new: test-up run-playwright-song-new test-down docker-clean-all
	@echo "✅ Song new tests finished. Report available in frontend/playwright-report/index.html"

# Lance uniquement les tests de la page d'ÉDITION d'une chanson
test-song-edit: test-up run-playwright-song-edit test-down docker-clean-all
	@echo "✅ Song edit tests finished. Report available in frontend/playwright-report/index.html"


# --- Test Helpers (ne pas lancer directement) ---

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
	@echo "--- Running ALL Playwright tests ---"
	@cd frontend && PLAYWRIGHT_TEST_BASE_URL=http://localhost:4001 npx playwright test

run-playwright-setlist:
	@echo "--- Running SETLIST Playwright tests ---"
	@cd frontend && PLAYWRIGHT_TEST_BASE_URL=http://localhost:4001 npx playwright test tests/setlist/

run-playwright-setlist-detail:
	@echo "--- Running SETLIST DETAIL Playwright tests ---"
	@cd frontend && PLAYWRIGHT_TEST_BASE_URL=http://localhost:4001 npx playwright test tests/setlist/detail.spec.ts

run-playwright-setlist-add:
	@echo "--- Running SETLIST ADD Playwright tests ---"
	@cd frontend && PLAYWRIGHT_TEST_BASE_URL=http://localhost:4001 npx playwright test tests/setlist/add.spec.ts

run-playwright-setlist-new:
	@echo "--- Running SETLIST NEW Playwright tests ---"
	@cd frontend && PLAYWRIGHT_TEST_BASE_URL=http://localhost:4001 npx playwright test tests/setlist/new.spec.ts

run-playwright-song:
	@echo "--- Running SONG Playwright tests ---"
	@cd frontend && PLAYWRIGHT_TEST_BASE_URL=http://localhost:4001 npx playwright test tests/song/

run-playwright-song-list:
	@echo "--- Running SONG LIST Playwright tests ---"
	@cd frontend && PLAYWRIGHT_TEST_BASE_URL=http://localhost:4001 npx playwright test tests/song/list.spec.ts

run-playwright-song-new:
	@echo "--- Running SONG NEW Playwright tests ---"
	@cd frontend && PLAYWRIGHT_TEST_BASE_URL=http://localhost:4001 npx playwright test tests/song/new.spec.ts

run-playwright-song-edit:
	@echo "--- Running SONG EDIT Playwright tests ---"
	@cd frontend && PLAYWRIGHT_TEST_BASE_URL=http://localhost:4001 npx playwright test tests/song/edit.spec.ts

test-down:
	@echo "--- Tearing down test environment ---"
	@docker-compose -f docker-compose.test.yml --env-file .env.test down -v --remove-orphans