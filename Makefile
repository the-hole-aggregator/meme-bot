include .env
export

# Defining variables for all scripts
SCRIPTS_DIR := scripts
GIT_HOOKS_INIT := $(SCRIPTS_DIR)/git_hooks_init.sh


# Tasks to run each script
hooks-init:
	sh $(GIT_HOOKS_INIT)
	
lint:
	golangci-lint run
	
dev-up:
	docker compose up -d
	sleep 3
	migrate -path migrations -database $(DATABASE_URL) up
	go run cmd/main.go
	
dev-down:
	migrate -path migrations -database $(DATABASE_URL) down
	docker compose down
	pkill -f "go run cmd/main.go"
	
remove-db:
	migrate -path migrations -database $(DATABASE_URL) down
	docker compose down -v
	
show-db:
	psql $(DATABASE_URL)
	
	
help:
	@echo " - hooks_init: init git hooks"
	@echo " - lint: run linter checks"
	@echo " - dev-up: start development environment"
	@echo " - dev-down: stop development environment"
