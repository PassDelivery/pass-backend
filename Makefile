include .env
export

export PROJECT_ROOT=${shell pwd}




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
	@if [ -z $("action") ]; then \
		echo "Отсутствует необходимый параметр action. Пример make migrate-action action=up"; \
		exit 1; \
	fi; \
	
	@docker compose run pass-postgres-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@pass-postgres:5432/${pPOSTGRES_DB}?sslmode=disable \
		"$(action)"

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down