include .env
export

export PROJECT_ROOT=$(shell pwd)

env-up:
	@docker compose up -d todoapp-postgres

env-down:
	@docker compose down todoapp-postgres

env-log:
	@docker compose logs -f todoapp-postgres

env-cleanup:
	@printf "Очистить все volume файлы? Опасность потери данных. [Y/n]: " ans; \
	read ans; \
	if [ "$$ans" = "Y" ]; then \
	  docker compose down todoapp-postgres && \
	  rm -rf out/pgdata && \
	  echo "Файлы очищены"; \
	else \
	  echo "Очистка окружения отменена"; \
  	fi

env-port-forward:
	@docker compose up -d port-forwarder

env-port-close:
	@docker compose down port-forwarder

migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Отсутсует необходимый параметр seq. Пример: make migrate-create seq=init" \
		exit 1; \
	fi; \
	docker compose run --rm todoapp-postgres-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-up:
	@$(MAKE) migrate-action action=up

migrate-down:
	@$(MAKE) migrate-action action=down


migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Отсутсует необходимый параметр action Пример: migrate-action action=up" \
		exit 1; \
	fi; \
	docker compose run --rm --use-aliases todoapp-postgres-migrate \
		-path /migrations \
		-database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@todoapp-postgres/${POSTGRES_DB}?sslmode=disable" \
		"$(action)"