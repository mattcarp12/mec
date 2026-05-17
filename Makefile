.PHONY: devup devdown devlogs api web setup

# ==============================================================================
# Dev Infra (Docker)
# ==============================================================================

devup:
	@echo "Starting infrastructure..."
	docker compose up -d

devdown:
	@echo "Stopping infrastructure..."
	docker compose down

devlogs:
	docker compose logs -f

# ==============================================================================
# Application Execution
# ==============================================================================

# Runs the Go backend API
api:
	@echo "Starting Go API..."
	cd backend && go run cmd/api/main.go

# Runs the Next.js frontend dev server
web:
	@echo "Starting Next.js frontend..."
	cd frontend && npm run dev

# Add to your existing Makefile variables at the top
DB_URL="postgres://devuser:devpassword@localhost:5432/ecommercedb?sslmode=disable"

# ==============================================================================
# Database & Migrations
# ==============================================================================

migrate-up:
	@echo "Running migrations up..."
	migrate -path backend/db/migrations -database $(DB_URL) up

migrate-down:
	@echo "Running migrations down..."
	migrate -path backend/db/migrations -database $(DB_URL) down 1

sqlc:
	@echo "Generating Go code with sqlc..."
	cd backend && sqlc generate

seed:
	@echo "Seeding database..."
	cd backend && go run cmd/seed/main.go