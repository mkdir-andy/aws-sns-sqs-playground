.PHONY: dev build clean republish

build:
	docker compose build

clean:
	docker compose down -v --remove-orphans

republish:
	docker compose restart publisher	

dev: build
	docker compose up --force-recreate --remove-orphans -d