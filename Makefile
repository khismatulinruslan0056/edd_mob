APP_NAME=server
BIN_DIR=./bin
BIN=$(BIN_DIR)/$(APP_NAME)

MIGRATIONS_DIR=./migrations

ENV_FILE=.env

MAIN=./cmd/$(APP_NAME)/main.go

.PHONY: all build run migrate-up migrate-down test fmt clean help

all: build

build:
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN) $(MAIN)

run: build
	$(BIN)

migrate-up:
	@echo "Applying migrations..."
	@bash -c 'set -a && source $(ENV_FILE) && migrate -path $(MIGRATIONS_DIR) -database "$$DATABASE_URL" up'

migrate-down:
	@echo "Reverting migrations..."
	@bash -c 'set -a && source $(ENV_FILE) && migrate -path $(MIGRATIONS_DIR) -database "$$DATABASE_URL" down'

fmt:
	go fmt ./...

clean:
	rm -rf $(BIN_DIR)

help:
	@echo "Makefile команды:"
	@echo "  build         - собрать бинарник"
	@echo "  run           - запустить приложение"
	@echo "  migrate-up    - применить миграции (golang-migrate up)"
	@echo "  migrate-down  - откатить миграции (golang-migrate down)"
	@echo "  fmt           - отформатировать код"
	@echo "  clean         - удалить бинарники"
