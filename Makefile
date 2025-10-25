# ========================
# VARIABLES
# ========================
APP_NAME := supportplatform
API_CMD := ./cmd/server
WORKER_CMD := ./cmd/worker
GO := go
PORT := 8080
MIGRATE_PATH := ./internal/db/migrations
DB_URL := postgres://postgres:psql1412@localhost:5432/support_platform_db?sslmode=disable&TimeZone=UTC
DB_URL_RENDER := ${DATABASE_URL_RENDER}

# ========================
# DEFAULT
# ========================
.PHONY: default
default: run

# ========================
# BUILD
# ========================
.PHONY: build
build:
	$(GO) build -o bin/$(APP_NAME) $(API_CMD)/main.go

# ========================
# RUN DEV SERVER
# ========================
.PHONY: run
run:
	@echo "Starting dev server on port $(PORT)..."
	$(GO) run $(API_CMD)/main.go

# ========================
# RUN WORKER
# ========================
.PHONY: worker
worker:
	@echo "Starting worker..."
	$(GO) run $(WORKER_CMD)/main.go

# ========================
# TEST
# ========================
.PHONY: test
test:
	$(GO) test -v ./...

# ========================
# CLEAN
# ========================
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/

# ========================
# Migration Commands
# ========================

# generate migration
migrate-create:
	@read -p "Enter migration name: " NAME; \
	UP_FILE=$(MIGRATE_PATH)/$$(date +%s)_$$NAME.up.sql; \
	DOWN_FILE=$(MIGRATE_PATH)/$$(date +%s)_$$NAME.down.sql; \
	echo "CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";" > $$UP_FILE; \
	echo "" >> $$UP_FILE; \
	echo "-- CREATE TABLE <table_name> ..." >> $$UP_FILE; \
	echo "CREATE TABLE <table_name> (" >> $$UP_FILE; \
	echo "    id UUID PRIMARY KEY DEFAULT uuid_generate_v4()," >> $$UP_FILE; \
	echo "    created_at TIMESTAMPTZ DEFAULT NOW()," >> $$UP_FILE; \
	echo "    updated_at TIMESTAMPTZ DEFAULT NOW()" >> $$UP_FILE; \
	echo ");" >> $$UP_FILE; \
	echo "DROP TABLE <table_name>;" > $$DOWN_FILE; \
	echo "Migration files created:" $$UP_FILE "and" $$DOWN_FILE

# generate alter table migration
migrate-alter:
	@read -p "Enter migration name: " NAME; \
	UP_FILE=$(MIGRATE_PATH)/$$(date +%s)_$$NAME.up.sql; \
	DOWN_FILE=$(MIGRATE_PATH)/$$(date +%s)_$$NAME.down.sql; \
	echo "-- ALTER TABLE <table_name> ADD COLUMN ..." > $$UP_FILE; \
	echo "ALTER TABLE <table_name>" >> $$UP_FILE; \
	echo "    ADD COLUMN <column_name> <type>;" >> $$UP_FILE; \
	echo "" >> $$UP_FILE; \
	echo "-- Remember to add more ALTER TABLE statements if needed" >> $$UP_FILE; \
	echo "-- ALTER TABLE <table_name> DROP COLUMN ..." > $$DOWN_FILE; \
	echo "ALTER TABLE <table_name>" >> $$DOWN_FILE; \
	echo "    DROP COLUMN IF EXISTS <column_name>;" >> $$DOWN_FILE; \
	echo "" >> $$DOWN_FILE; \
	echo "Migration files created:" $$UP_FILE "and" $$DOWN_FILE


# run migration up
migrate-up:
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" up

# rollback last migration
migrate-down:
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" down 1

# check migration status (version + dirty flag)
migrate-status:
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" version || true

# force clean a dirty migration (set to a specific version)
migrate-force:
	@read -p "Enter version to force: " VERSION; \
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" force $$VERSION
	echo "Forced to version $$VERSION."

# ========================
# DEPLOY (for Render)
# ========================
.PHONY: deploy
deploy:
	@echo ">>> Installing migrate CLI if not installed..."
	go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest

	@echo ">>> Building binary..."
	make build

	@echo ">>> Running migrations..."
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL_RENDER)" up

	@echo ">>> Starting server..."
	./bin/$(APP_NAME)
