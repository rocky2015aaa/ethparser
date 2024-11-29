api-docs: ## Generate API docs with swaggo
	@echo "========================="
	@echo "Generate Swagger API Docs"
	@echo "========================="
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init swag init --parseDependency --parseInternal

test: ## Run unit tests
test:
	@echo "=================="
	@echo "Running unit tests"
	@echo "=================="
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

run: ## Run the server
run: 
	go run main.go