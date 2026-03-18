# === DEFAULT ===
default:
	just --list

# === ALL ===
dev:
	just backend-run &
	just frontend-dev

# === BACKEND ===
backend-run:
	cd backend && go run cmd/api/main.go

backend-smoke-endpoints:
	powershell -ExecutionPolicy Bypass -File scripts/smoke-endpoints.ps1

backend-lint:
	cd backend && golangci-lint run

backend-fmt:
	cd backend && gofmt -s -w .
	cd backend && goimports -w .

backend-tidy:
	cd backend && go mod tidy
	cd backend && go mod verify

# === FRONTEND ===
frontend-dev:
	cd frontend && npm run dev

frontend-format:
	cd frontend && npm run format

frontend-lint:
	cd frontend && npm run lint

frontend-check:
	just frontend-lint
	just frontend-format

# === DOCS ===
backend-swagger:
	go install github.com/swaggo/swag/cmd/swag@latest
	cd backend && swag init -g cmd/api/main.go --parseDependency --parseInternal

# === SETUP ===
setup:
	just install-tools
	cd frontend && npm install
	cd backend && go mod tidy
	cd backend && go mod verify

# === TOOLS ===
install-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest
	-pre-commit install
	-pre-commit install --hook-type commit-msg
	cd frontend && npm install --save-dev prettier eslint
