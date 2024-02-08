DATABASE_PATH = ./db/database.db

init-db:
	@echo "Initializing the database"
	@$(MAKE) migrate-up

migrate-up:
	@echo "Running database up-migrations"
	@./build/go-webhook-server migrate-up --database $(DATABASE_PATH)

migrate-down:
	@echo "Rolling back database via down-migrations"
	@./build/go-webhook-server migrate-down --database $(DATABASE_PATH)

build:
	@echo "Building the application"
	@go build -o "./build/go-webhook-server"