# Set the service name
SERVICE_NAME = user-api

# Docker Image Name
IMAGE_NAME = user-api

# 🛠️ Build the Golang Docker Image
build:
	docker build -t $(IMAGE_NAME) .

# 🚀 Run the Application with Docker
run:
	docker run -p 8080:8080 --env-file .env $(IMAGE_NAME)

# 🏗️ Build and Start All Services (Postgres, RabbitMQ, API)
start:
	docker-compose up --build -d

# ⏹️ Stop and Remove Containers
stop:
	docker-compose down

# 📜 View Running Containers
ps:
	docker ps

# 🔍 View Application Logs
logs:
	docker-compose logs -f

# 🧪 Run Unit Tests Inside Container
test-unit:
	docker-compose run --rm test go test -tags=unit_test ./pkg/...

# 🧪 Run Integration Tests Inside Container
test-integration:
	docker-compose run --rm test go test -tags=integration ./pkg/repository/

# 🏁 Run All Tests Inside Container
test:
	docker-compose run --rm test go test ./...

# 📜 Clean Docker Resources (Stopped Containers, Networks, Images)
clean:
	docker system prune -f

# 🆘 Help Menu
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@echo "  build             - Build the Docker image"
	@echo "  run               - Run the API container"
	@echo "  start             - Start all services (API, Postgres, RabbitMQ)"
	@echo "  stop              - Stop all running services"
	@echo "  ps                - Show running containers"
	@echo "  logs              - View logs"
	@echo "  test-unit         - Run unit tests inside the container"
	@echo "  test-integration  - Run integration tests inside the container"
	@echo "  test              - Run all tests inside the container"
	@echo "  clean             - Remove unused Docker resources"
