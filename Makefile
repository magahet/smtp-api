.PHONY: build run update start stop remove restart
	
build:
	docker build -t penpal-api .

run:
	docker run --link mongo --rm --name penpal-api -d -p 8001:8001 --env-file env_vars penpal-api

update: stop build run

start:
	docker start penpal-api

stop:
	docker stop penpal-api

restart: stop run

db-run:
	docker run --rm --name penpal-db -d -p 27017:27017 mongo

db-start:
	docker start penpal-db

db-stop:
	docker stop penpal-db

db-ui-run:
	docker run -d --rm --name mongo-express --link mongo:mongo -p 8111:8081 mongo-express