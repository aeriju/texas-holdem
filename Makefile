# Texas Hold'em - Build and deployment automation
.PHONY: backend frontend docker-build docker-backend docker-frontend test loadtest clean

# Default target
all: backend frontend

# --- Backend (Go) ---
backend:
	@echo "Building backend..."
	cd backend && go build -o bin/server ./cmd/api
	@echo "Backend built: backend/bin/server"

backend-run: backend
	cd backend && ./bin/server

# --- Frontend (Flutter) ---
frontend:
	@echo "Building Flutter web app..."
	cd frontend && flutter pub get && flutter build web --no-wasm-dry-run
	@echo "Frontend built: frontend/build/web"

# --- Docker builds (AMD64 for GKE) ---
docker-build:
	@echo "Building Docker images for linux/amd64..."
	docker build --platform linux/amd64 -t holdem-backend:latest -f backend/Dockerfile .
	docker build --platform linux/amd64 -t holdem-frontend:latest -f frontend/Dockerfile frontend
	@echo "Docker images built."

docker-backend:
	@echo "Building backend Docker image for linux/amd64..."
	docker build --platform linux/amd64 -t holdem-backend:latest -f backend/Dockerfile .
	@echo "Backend Docker image built."

docker-frontend:
	@echo "Building frontend Docker image for linux/amd64..."
	docker build --platform linux/amd64 -t holdem-frontend:latest -f frontend/Dockerfile frontend
	@echo "Frontend Docker image built."

# --- Tests ---
test:
	cd backend && go test -v ./...

# --- Load testing ---
loadtest:
	@echo "Running load tests..."
	cd loadtest && k6 run loadtest.js
	@echo "Load test complete."

# --- Clean ---
clean:
	rm -rf backend/bin
	rm -rf frontend/build
