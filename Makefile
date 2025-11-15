NETWORK=core-net

.PHONY: up-all up-postgres migrate-user migrate-robokassa up-services down build clean logs-app-docker pre-commit-init create-network

pre-commit-init:
	pre-commit install

# Создать сеть, если её нет
create-network:
	docker network create $(NETWORK) || true

# Поднять только Postgres
up-postgres: create-network
	docker compose up -d postgres

# Миграции user_service
migrate-user:
	docker compose run --rm goose

# Миграции robokassa_service
migrate-robokassa:
	docker compose run --rm goose_robokassa

# Поднять сервисы после миграций
up-services:
	docker compose up -d user_service robokassa_service

# Полный апстрим: создать сеть → Postgres → миграции → сервисы
up-all: up-postgres migrate-user migrate-robokassa up-services

# Остановить все сервисы
down:
	docker compose down

# Просмотр логов
logs-app-docker:
	docker compose logs -f

# Собрать образы
build:
	docker compose build

# Полная очистка
clean:
	docker compose down -v
	docker network rm $(NETWORK) || true
