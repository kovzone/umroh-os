# SERVICES lists every scaffolded service under services/.
# gateway-svc and broker-svc are deferred (gateway: next session; broker: ADR 0006).
SERVICES = services/gateway-svc services/iam-svc services/catalog-svc services/booking-svc services/jamaah-svc services/payment-svc services/visa-svc services/ops-svc services/logistics-svc services/finance-svc services/crm-svc

# Local migration URL — single shared database `umrohos_dev` per ADR 0007.
# Host port 5432 is mapped from the postgres container (see docker-compose.dev.yml).
LOCAL_DB_URL = "postgres://postgres:changeme@localhost:5432/umrohos_dev?sslmode=disable"

help: ## Show available commands
	@echo ""
	@echo "Usage: make <target>"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'
	@echo ""

# ===========================================
# Code Generation
# ===========================================

sqlc: ## Generate store code (sqlc) for all services
	@for svc in $(SERVICES); do \
		if [ -f $$svc/store/postgres_store/sqlc.yaml ]; then \
			echo "sqlc generate: $$svc"; \
			cd $$svc/store/postgres_store && sqlc generate && cd - > /dev/null; \
		fi; \
	done

oapi: ## Generate REST API code (oapi-codegen) for all services
	@for svc in $(SERVICES); do \
		if [ -f $$svc/api/rest_oapi/.oapi-codegen.yaml ]; then \
			echo "oapi-codegen: $$svc"; \
			cd $$svc/api/rest_oapi && oapi-codegen --config=.oapi-codegen.yaml openapi.yaml && cd - > /dev/null; \
		fi; \
	done

genpb: ## Generate protobuf code for all services
	@for svc in $(SERVICES); do \
		if ls $$svc/api/grpc_api/pb/*.proto >/dev/null 2>&1; then \
			echo "protoc: $$svc"; \
			cd $$svc && protoc --proto_path=api/grpc_api/pb api/grpc_api/pb/*.proto \
				--go_out=api/grpc_api/pb --go_opt=paths=source_relative \
				--go-grpc_out=api/grpc_api/pb --go-grpc_opt=paths=source_relative && cd - > /dev/null; \
		fi; \
	done

generate: sqlc oapi genpb ## Run all code generation (sqlc + oapi-codegen + protoc)

# ===========================================
# Database Migrations (golang-migrate, per ADR 0007)
# ===========================================

migrate-up: ## Apply all pending migrations
	migrate -source file://migration -database $(LOCAL_DB_URL) up

migrate-down: ## Roll back migrations (usage: make migrate-down STEPS=1)
	migrate -source file://migration -database $(LOCAL_DB_URL) down $(STEPS)

migrate-version: ## Show current migration version
	migrate -source file://migration -database $(LOCAL_DB_URL) version

migrate-force: ## Force migration version (usage: make migrate-force VERSION=1)
	migrate -source file://migration -database $(LOCAL_DB_URL) force $(VERSION)

migrate-create: ## Create new migration pair (usage: make migrate-create NAME=create_iam_users)
	migrate create -ext sql -dir migration -seq $(NAME)

# ===========================================
# Docker Compose
# ===========================================

dev-up: ## Start the dev environment (docker compose)
	docker compose -f docker-compose.dev.yml up -d

dev-down: ## Stop the dev environment
	docker compose -f docker-compose.dev.yml down

dev-down-v: ## Stop the dev environment and remove volumes
	docker compose -f docker-compose.dev.yml down -v

dev-logs: ## Tail logs for all services
	docker compose -f docker-compose.dev.yml logs -f

dev-ps: ## Show running containers
	docker compose -f docker-compose.dev.yml ps

dev-bootstrap: dev-up ## Start dev environment then apply all migrations
	@echo "Waiting for postgres to accept connections..."
	@for i in 1 2 3 4 5 6 7 8 9 10; do \
		docker compose -f docker-compose.dev.yml exec -T postgres pg_isready -U postgres -d umrohos_dev > /dev/null 2>&1 && break; \
		echo "  ...postgres not ready yet (attempt $$i/10), waiting 2s"; \
		sleep 2; \
	done
	@$(MAKE) migrate-up

# ===========================================
# Testing
# ===========================================

test: ## Run unit tests for all services
	@for svc in $(SERVICES); do \
		if [ -f $$svc/go.mod ]; then \
			echo "Testing: $$svc"; \
			cd $$svc && go test ./... -cover && cd - > /dev/null; \
		fi; \
	done

test-v: ## Run unit tests verbose for all services
	@for svc in $(SERVICES); do \
		if [ -f $$svc/go.mod ]; then \
			echo "Testing: $$svc"; \
			cd $$svc && go test ./... -v -cover && cd - > /dev/null; \
		fi; \
	done

test-svc: ## Run tests for a specific service (usage: make test-svc SVC=services/iam-svc)
	cd $(SVC) && go test ./... -v -cover

test-coverage: ## Generate coverage report for a specific service (usage: make test-coverage SVC=services/iam-svc)
	cd $(SVC) && go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: $(SVC)/coverage.html"

# ===========================================
# Docker Image Management
# ===========================================

dev-rm-all: ## Remove all service Docker images (idempotent; no-ops on missing)
	@for svc in $(SERVICES); do \
		img=$$(basename $$svc); \
		if docker image inspect $$img >/dev/null 2>&1; then \
			echo "Removing $$img"; \
			docker rmi $$img >/dev/null; \
		else \
			echo "Skipping $$img (no image)"; \
		fi; \
	done

dev-rebuild: ## Rebuild and restart a specific service (usage: make dev-rebuild SVC=iam-svc)
	docker compose -f docker-compose.dev.yml up -d --build $(SVC)

# ===========================================
# E2E Testing (Playwright, per ADR 0008)
# ===========================================

e2e-install: ## Install e2e dependencies (one-time / CI)
	cd tests/e2e && npm install

e2e: ## Run the full e2e suite against the running stack
	cd tests/e2e && npm test

.PHONY: help sqlc oapi genpb generate migrate-up migrate-down migrate-version migrate-force migrate-create dev-up dev-down dev-down-v dev-logs dev-ps dev-bootstrap test test-v test-svc test-coverage dev-rm-all dev-rebuild e2e-install e2e
