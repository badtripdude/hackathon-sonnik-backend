NETWORK=core-net

.PHONY: up down build clean create-network pre-commit-init

pre-commit-init:
	pre-commit install

create-network:
	docker network create $(NETWORK) || true

up: create-network
#	docker compose -f another.yaml up -d
	docker compose up -d

down:
	docker compose down
#	docker compose -f another.yaml down

logs-app-docker:
	docker compose logs -f

build:
	docker compose build
#	docker compose -f another.yaml build

clean:
	docker compose down -v
	docker network rm core-net || true
