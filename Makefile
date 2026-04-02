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
	

build:
	GOOS=linux GOARCH=amd64 go build -o bot ./cmd
	
dev-up:
	docker compose up -d postgres
	sleep 5
	docker compose run --rm migrate
	docker compose up -d --build bot
	
dev-down:
	docker compose down
	
remove-db:
	docker compose run --rm migrate -path migrations -database "$DATABASE_URL" down
	docker compose down -v
	
show-db:
	docker exec -it memes_db psql "$(DATABASE_URL)"
	
logs:
	docker compose logs -f bot
	
test:
	@echo "Tests running..."
	go test -v ./...
	@echo "Tests completed"
	
	
help:
	@echo " - hooks_init: init git hooks"
	@echo " - lint: run linter checks"
	@echo " - dev-up: start development environment"
	@echo " - dev-down: stop development environment"
