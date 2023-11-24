# To use env variables from local .env file you need to install
# npm install -g dotenv-cli  
include .env
export

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: run
run: ## Runs the whole app in docker containers ( --abort-on-container-exit ) 
	docker compose -f ./docker-compose.yaml up

.PHONY: build
build: ## Runs the whole app in docker containers ( --abort-on-container-exit ) 
	docker compose -f ./docker-compose.yaml build

.PHONY: down
down: ## Stops containers, networks, volumes, and images created by up
	docker compose -f ./docker-compose.yaml down

.PHONY: migrate postgresql up
migrateup: ## Stops containers, networks, volumes, and images created by up
	migrate -path ./db/migrations -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" up

.PHONY: migrate postgresql down
migratedown: ## Stops containers, networks, volumes, and images created by up
	migrate -path ./db/migrations -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" down