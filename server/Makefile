# Makefile

.PHONY: build run

build:
	docker-compose build

run: build
	docker-compose up

stop:
	docker-compose down

logs:
	docker-compose logs -f
