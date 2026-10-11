include .env
export

export PROJECT_ROOT=${shell pwd}


env-up:
	@docker compose up -d pass-postgres

env-down:
	@docker compose down pass-postgres


env-cleanup:
	@read -p "Очитсить все volume файлы окружения? Данные будут утеряны. [y/n]: " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down pass-postgres port-forwarder && \
		rm -rf out/pgdata && \
		echo "Файлы окружения очищены."; \
	else \
		echo "Очистка окружения отменена."; \
	fi


migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Отсутствует параметр seq. Пример make migrate-create seq=init"; \
		exit 1; \
	fi; \

	docker compose run --rm pass-postgres-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"


migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Отсутствует необходимый параметр action. Пример make migrate-action action=up"; \
		exit 1; \
	fi; \
	
	@docker compose run --rm pass-postgres-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@pass-postgres:5432/${POSTGRES_DB}?sslmode=disable \
		"$(action)"

migrate-up:
	@make migrate-action action=up
	
migrate-down:
	@make migrate-action action=down

forwarder-up:
	@docker compose up -d port-forwarder

forwarder-down:
	@docker compose down port-forwarder



pass-run:
	@export LOGGER_PATH=${PROJECT_ROOT}/out/logs && \
	export POSTGRES_HOST=localhost && \
	go mod tidy && \
	go run cmd/passapp/main.go