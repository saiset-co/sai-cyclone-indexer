SERVICE_NAME=sai-cyclone-indexer
EXTERNAL_PORT=8881
PORT=8881

build:
	docker-compose up -d --build

up:
	docker-compose up -d

sh:
	docker-compose exec ${SERVICE_NAME} bash

log:
	docker-compose logs -f ${SERVICE_NAME}

down:
	docker-compose down

test:
	go test ./tests -run TestStart -count=1

docker:
	docker build -t ${SERVICE_NAME} .
	docker stop ${SERVICE_NAME} || true
	docker rm ${SERVICE_NAME} || true
	docker run -d -p ${EXTERNAL_PORT}:${PORT} --restart unless-stopped --name ${SERVICE_NAME} ${SERVICE_NAME}
