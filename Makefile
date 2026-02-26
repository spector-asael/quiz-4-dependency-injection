# CMPS3162 Advanced Databases - Banking System
# Asael Tobar
# February 23rd, 2025

# Pass in the .envrc file, which exports BANK_DB_DSN
include .envrc

## run: run the cmd/api application
.PHONY: run help checkbalance deposit comment healthcheck all
run: 
	@echo 'Running application...'
	@go run ./cmd/api -db-dsn="${BANK_DB_DSN}"

# Help target
help:
	@echo ""
	@echo "Application:"
	@echo "  make run            - Run API server"
	@echo ""
	@echo "API Testing (requires server running):"
	@echo "  make checkbalance   - Test GET /v1/balance (user 1)"
	@echo "  make checkbalance2  - Test POST /v1/balance (user 2)"
	@echo "  make deposit        - Test POST /v1/deposit (valid request)"
	@echo "  make deposit2       - Test POST /v1/deposit (invalid request)"
	@echo ""
	@echo "Database:"
	@echo "  make db/psql        - Connect to database using psql"
	@echo ""
	@echo "Migrations:"
	@echo "  make db/migrations/new name=NAME"
	@echo "                      - Create new migration"
	@echo ""
	@echo "  make migrations/up  - Apply all migrations"
	@echo "  make migrations/down"
	@echo "                      - Revert all migrations"
	@echo ""
	@echo "  make migrations/fix version=VERSION"
	@echo "                      - Force schema_migrations to VERSION"
	@echo ""
	@echo "  make migrations/init"
	@echo "                      - Generate all initial banking tables"
	@echo ""
	@echo "Environment:"
	@echo "  Requires BANK_DB_DSN exported in .envrc"
	@echo ""

# POST /balance
checkbalance:
	@echo "Testing /balance..."
	curl -X GET http://localhost:4000/v1/balance \
	-H "Content-Type: application/json" \
	-d '{"gl_account_id":1}'

checkbalance2:
	@echo "Testing /balance..."
	curl -X POST http://localhost:4000/v1/balance \
	-H "Content-Type: application/json" \
	-d '{"user_id":2,"bank_number":111111}'

# POST /deposit
deposit:
	@echo "Testing /deposit..."
	curl -X POST http://localhost:4000/v1/deposit \
	-H "Content-Type: application/json" \
	-d '{"gl_account_id":1,"amount":500.75}'

deposit2:
	@echo "Testing /deposit..."
	curl -X POST http://localhost:4000/v1/deposit \
	-H "Content-Type: application/json" \
	-d '{"user_id":1,"bank_number":111111,"deposit_amount":500 75}'

checkhistory:
	@echo "Testing /history..."
	curl -X POST http://localhost:4000/v1/history \
	-H "Content-Type: application/json" \
	-d '{"user_id":1}'

transfer:
	@echo "Testing /transfer..."
	curl -X POST http://localhost:4000/v1/transfer \
	-H "Content-Type: application/json" \
	-d '{ "from_account_id": 1, "to_account_id": 2, "amount": 100.00}'

deleteentry:
	@echo "Testing /delete..."
		curl -X DELETE http://localhost:4000/v1/delete \
		-H "Content-Type: application/json" \
		-d '{"ledger_id":3}'

updateentry:
	@echo "Testing /update..."
		curl -X PATCH http://localhost:4000/v1/update \
		-H "Content-Type: application/json" \
		-d '{"ledger_id":4,"amount":1000.50}'

## db/psql: Connect to the banking database using psql
.PHONY: db
db:
	psql ${BANK_DB_DSN}

## db/migrations/new name=$1: Create a new database migration
.PHONY: migrations/new
migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## migrations/up: Apply all up database migrations
.PHONY: migrations/up
migrations/up:
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${BANK_DB_DSN} up

## db/migrations/down: Revert all migrations
.PHONY: migrations/down
migrations/down:
	@echo 'Reverting all migrations...'
	migrate -path ./migrations -database ${BANK_DB_DSN} down

## db/migrations/fix version=$1: Force schema_migrations version
.PHONY: migrations/fix
migrations/fix:
	@echo 'Forcing schema migrations version to ${version}...'
	migrate -path ./migrations -database ${BANK_DB_DSN} force ${version}

