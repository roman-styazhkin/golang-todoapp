include .env
export

export PROJECT_ROOT=$(shell pwd)

env-up:
	@docker compose up -d todoapp-postgres

env-down:
	@docker compose down todoapp-postgres

env-cleanup:
	@read -p "Очистить все volume файлы окружения? Опасность утери данных. [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down todoapp-postgres && \
		rm -rf ${PROJECT_ROOT}/out/pgdata && \
		echo "Файлы окружения очищены."; \
	else \
	  echo "Отмена очистки файлов окружения"; \
	fi;

migrate-create:
	@if [ -z "$(seqName)" ]; then \
		echo "Не задан обязательный параметр seqName. Пример команды: make migrate-create seqName=init"; \
		exit 1; \
  	fi;

	@docker compose run --rm todoapp-postgres-migrate \
	create -ext sql -dir /migrations -seq $(seqName)

migrate-up:
	@if [ -z "$(steps)" ]; then \
		echo "Не задан обязательный параметр steps. Пример команды: make migrate-up steps=1"; \
		exit 1; \
	fi;

	@docker compose run --rm todoapp-postgres-migrate \
	-path /migrations -database "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@todoapp-postgres:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable" \
	 up $(steps)

migrate-down:
	@if [ -z "$(steps)" ]; then \
		echo "Не задан обязательный параметр steps. Пример команды: make migrate-down steps=1"; \
		exit 1; \
	fi;

	@docker compose run --rm todoapp-postgres-migrate \
	-path /migrations -database "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@todoapp-postgres:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable" \
	 down $(steps)


port-forwarder:
	@docker compose up -d postgres-forwarder

todoapp-run:
	@export LOGGER_FOLDER=${PROJECT_ROOT}/out/logs && \
	export POSTGRES_HOST=localhost && \
	go run cmd/api/main.go
