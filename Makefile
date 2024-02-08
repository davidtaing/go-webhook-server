DATABASE_PATH = db/database.db
DATABASE_URL = sqlite3://${DATABASE_PATH}

init-db:
	@echo "Deleting database at ${DATABASE_PATH}"
	@rm -f ${DATABASE_PATH}
	@echo "Creating empty database at ${DATABASE_PATH}"
	@sqlite3 ${DATABASE_PATH} "VACUUM;"
	@$(MAKE) migrate-up

migrate-up:
	@echo "Running database up-migrations"
	@migrate -database $(DATABASE_URL) -path db/migrations up

migrate-down:
	@echo "Rolling back database via down-migrations"
	@migrate -database $(DATABASE_URL) -path db/migrations down
