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
	  docker compose down todoapp-postgres port-forwarder && \
	  rm -rf ${PROJECT_ROOT}/out/pgdata && \
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

logs-cleanup: ## env: Очистить файлы логов из out/logs
	@read -p "Очистить все log файлы? Опасность утери логов. [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		rm -rf ${PROJECT_ROOT}/out/logs && \
		echo "Файлы логов очищены"; \
	else \
		echo "Очистка логов отменена"; \
	fi

todoapp-deploy:
	docker compose up -d --build todoapp

todoapp-run:
	@export LOGGER_FOLDER=${PROJECT_ROOT}/out/logs && \
	export POSTGRES_HOST=localhost && \
	go mod tidy && \
	go run ${PROJECT_ROOT}/cmd/todoapp/main.go

ps:
	@docker compose ps


#docker logs todoapp