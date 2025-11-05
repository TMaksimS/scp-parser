.PHONY: swag test_db run_parser down_local

swag_generate:
	go run github.com/swaggo/swag/cmd/swag@latest init -d ./server/cmd/api,./server/internal/handlers,./server/internal/service,./server/domain -o ./docs

test_db:
	@docker compose -f deployment/local/docker-compose-local.yaml --env-file .env up -d
run_parser: test_db
	@echo "CMD start for parser.go has been running" 
	@go run -race parser/cmd/parser/main.go

run_server: test_db swag_generate
	@echo "CMD start for server has been running"
	@go run -race server/cmd/api/main.go

down_local:
	@docker compose -f deployment/local/docker-compose-local.yaml --env-file .env down --remove-orphans