# To use env variables from local .env file you need to install
# npm install -g dotenv-cli  
include .env
export

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: run
run: ## Runs the whole app in docker containers 
	docker compose -f ./docker-compose.yaml up -d 

.PHONY: build
build: ## Build or rebuild containers
	docker compose -f ./docker-compose.yaml build

.PHONY: down
down: ## Remove docker containers
	docker compose -f ./docker-compose.yaml down

.PHONY: migrate postgresql up
migrateup: ## migrate postgresql up
	migrate -path ./db/migrations -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" up

.PHONY: migrate postgresql down
migratedown: ## migrate postgresql down
	migrate -path ./db/migrations -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" down

.PHONY: cert
cert:
	mkdir -p cert
	openssl genrsa -out cert/access 4096
	openssl rsa -in cert/access -pubout -out cert/access.pub

.PHONY: goModule
goModule: ## Remove docker containers
	go env -w GO111MODULE=on
