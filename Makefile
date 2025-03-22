# Define variables
APP_NAME = app
FRONTEND_DIR = frontend
BACKEND_DIR = backend
DIST_DIR = $(BACKEND_DIR)/dist

.DEFAULT_GOAL := help

# Show help
help:
	@echo "Available commands:"
	@echo "  make help          - Show this help message"
	@echo "  make build         - Build the entire application (frontend + backend)"
	@echo "  make build-frontend- Build only the frontend"
	@echo "  make build-backend - Build only the backend (requires frontend build)"
	@echo "  make run          - Build and run the application"
	@echo "  make run-frontend  - Run the frontend development server"
	@echo "  make run-backend   - Run the backend development server"
	@echo "  make run-dev      - Run both frontend and backend development servers"
	@echo "  make clean        - Clean up generated files"
	@echo "  make rebuild      - Clean and rebuild the entire application"

# Default target
all: build

# Always run the build
build: build-backend

# Build frontend
build-frontend:
	@echo "Building frontend..."
	cd $(FRONTEND_DIR) && npm install && npm run build
	@echo "Moving frontend dist to backend..."
	rm -rf $(DIST_DIR)  # Remove old files
	mv $(FRONTEND_DIR)/dist $(DIST_DIR)  # Move built frontend to backend

# Build backend
build-backend: build-frontend
	@echo "Building Go binary..."
	cd $(BACKEND_DIR) && CGO_ENABLED=1 go build -o ../$(APP_NAME) ./...

# Run the application
run: build
	@echo "Starting the app..."
	./$(APP_NAME)

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
	rm -rf $(APP_NAME) $(DIST_DIR)

# Rebuild everything
rebuild: clean build

# Ensure these targets are always executed
.PHONY: all build build-frontend build-backend run run-frontend run-backend run-dev clean rebuild