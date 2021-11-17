default: build

build: docker

.PHONY: docker
docker:
	docker-compose up

.PHONY: rabbitmq
rabbitmq:
	-docker-compose up -d rabbitmq

.PHONY: aws
aws:
	docker-compose up -d aws

.PHONY: nats
nats:
	docker-compose up -d nats

.PHONY: gcp
gcp:
	docker-compose up -d pubsub

.PHONY: pubsub
pubsub: gcp

.PHONY: clean
clean:
	-docker-compose down -d rabbitmq
