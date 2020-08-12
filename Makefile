.PHONY: build run create start stop remove

build:
	docker build -t smtp-api .

run:
	docker run --rm --name smtp-api -d -p 8083:80 --env-file env_vars smtp-api

start:
	docker start smtp-api

stop:
	docker stop smtp-api

remove:
	docker rm smtp-api
