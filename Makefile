-include .env
export

up:
	docker compose up -d --build backend frontend db

down:
	docker compose down

logs:
	docker compose logs -f backend frontend db

migrate:
	docker compose up migrator

migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Erreur : L'argument 'name' est requis. Usage : make migrate-create name=mon_nom_de_migration"; \
		exit 1; \
	fi
	cd backend && migrate create -ext sql -dir db/migrations -seq $(name)

deploy: migrate up


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

test-all: test-unit test
	@echo "✅ All tests (Unit & E2E) finished successfully."


test-unit:
	@echo "--- Running Unit Tests ---"
	@cd frontend && npx vitest run
	@echo "✅ Unit Tests finished."

test-unit-watch:
	@echo "--- Running Unit Tests in watch mode ---"
	@cd frontend && npx vitest


# --- End-to-End Tests (Playwright) ---
test: test-up run-playwright test-down docker-clean-all
	@echo "✅ E2E Tests finished. Report available in frontend/playwright-report/index.html"


test-setlist: test-up run-playwright-setlist test-down docker-clean-project
	@echo "✅ Setlist tests finished."

test-setlist-detail: test-up run-playwright-setlist-detail test-down docker-clean-project
	@echo "✅ Setlist detail tests finished."

test-setlist-add: test-up run-playwright-setlist-add test-down docker-clean-project
	@echo "✅ Setlist add tests finished."

test-setlist-new: test-up run-playwright-setlist-new test-down docker-clean-project
	@echo "✅ Setlist new tests finished."

test-setlist-duplicate: test-up run-playwright-setlist-duplicate test-down docker-clean-project
	@echo "✅ Setlist duplicate test finished."

test-song: test-up run-playwright-song test-down docker-clean-project
	@echo "✅ Song tests finished."

test-song-list: test-up run-playwright-song-list test-down docker-clean-project
	@echo "✅ Song list tests finished."

test-song-new: test-up run-playwright-song-new test-down docker-clean-project
	@echo "✅ Song new tests finished."

test-song-edit: test-up run-playwright-song-edit test-down docker-clean-project
	@echo "✅ Song edit tests finished."

test-multi-group: test-up run-playwright-multi-group test-down docker-clean-project
	@echo "✅ Multi-group E2E test finished."

test-settings: test-up run-playwright-settings test-down docker-clean-project
	@echo "✅ Settings tests finished."

test-settings-account: test-up run-playwright-settings-account test-down docker-clean-project
	@echo "✅ Settings account test finished."

test-settings-members: test-up run-playwright-settings-members test-down docker-clean-project
	@echo "✅ Settings members test finished."

test-interlude: test-up run-playwright-interlude test-down docker-clean-project
	@echo "✅ Interlude tests finished."

test-interlude-new: test-up run-playwright-interlude-new test-down docker-clean-project
	@echo "✅ Interlude new tests finished."

test-interlude-behavior: test-up run-playwright-interlude-behavior test-down docker-clean-project
	@echo "✅ Interlude behavior tests finished."

test-up:
	@echo "--- Cleaning up previous test environment ---"
	@docker compose -f docker-compose.test.yml --env-file .env.test down -v --remove-orphans
	@echo "--- Building and starting test environment (DB, Backend with seed, Frontend) ---"
	@docker compose -f docker-compose.test.yml --env-file .env.test up --build -d || \
        (echo "🔴 'docker compose up' failed. Displaying logs:"; \
        docker compose -f docker-compose.test.yml --env-file .env.test logs; \
        exit 1)
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

run-playwright-setlist-duplicate:
	@echo "--- Running SETLIST DUPLICATE Playwright test ---"
	@cd frontend && npx playwright test tests/setlist/duplicate.spec.ts

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

run-playwright-multi-group:
	@echo "--- Running MULTI-GROUP Playwright test ---"
	@cd frontend && npx playwright test tests/group/multi-group.spec.ts

run-playwright-settings:
	@echo "--- Running SETTINGS Playwright tests ---"
	@cd frontend && npx playwright test tests/settings/

run-playwright-settings-account:
	@echo "--- Running SETTINGS ACCOUNT Playwright test ---"
	@cd frontend && npx playwright test tests/settings/account.spec.ts

run-playwright-settings-members:
	@echo "--- Running SETTINGS MEMBERS Playwright test ---"
	@cd frontend && npx playwright test tests/settings/members.spec.ts

run-playwright-interlude:
	@echo "--- Running INTERLUDE Playwright tests ---"
	@cd frontend && npx playwright test tests/interlude/

run-playwright-interlude-new:
	@echo "--- Running INTERLUDE NEW Playwright tests ---"
	@cd frontend && npx playwright test tests/interlude/new.spec.ts

run-playwright-interlude-behavior:
	@echo "--- Running INTERLUDE BEHAVIOR Playwright tests ---"
	@cd frontend && npx playwright test tests/interlude/behavior.spec.ts

test-down:
	@echo "--- Tearing down test environment ---"
	@docker compose -f docker-compose.test.yml --env-file .env.test down -v --remove-orphans