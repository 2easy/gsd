# Define variables
FRONTEND_DIR = frontend
BACKEND_DIR = backend
DIST_DIR = $(BACKEND_DIR)/dist

.DEFAULT_GOAL := help

# Show help
help:
	@echo "Available commands:"
	@echo "  make help          - Show this help message"
	@echo "  make build-frontend- Build the frontend"
	@echo "  make docker-build  - Build the Docker image"
	@echo "  make run-frontend  - Run the frontend development server"
	@echo "  make run-backend   - Run the backend development server"
	@echo "  make run-dev      - Run both frontend and backend development servers"
	@echo "  make clean        - Clean up generated files"

# Build frontend
build-frontend:
	@echo "Building frontend..."
	cd $(FRONTEND_DIR) && npm install && npm run build
	@echo "Moving frontend dist to backend..."
	rm -rf $(DIST_DIR)  # Remove old files
	mv $(FRONTEND_DIR)/dist $(DIST_DIR)  # Move built frontend to backend

# Build Docker image
docker-build: build-frontend
	@echo "Building Docker image..."
	docker build -t gsd .

# Run the frontend development server
run-frontend:
	@echo "Starting the frontend development server..."
	cd $(FRONTEND_DIR) && npm run dev

# Run the backend development server
run-backend:
	@echo "Starting the backend Go server..."
	cd $(BACKEND_DIR) && go run .

# Run the development servers
run-dev:
	@echo "Starting the frontend development server..."
	cd $(FRONTEND_DIR) && npm run dev &
	@echo "Starting the backend Go server..."
	cd $(BACKEND_DIR) && go run . &
	@echo "Both servers are running."

# Clean up generated files
clean:
	@echo "Cleaning up..."
	rm -rf $(DIST_DIR)

# Ensure these targets are always executed
.PHONY: help build-frontend docker-build run-frontend run-backend run-dev clean