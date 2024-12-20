# AnhCao 2024
DOCKER_USERNAME = anhcaoo
DOCKER_IMAGE = electric-notifications
TAGGED_VERSION = 1.0.0
DOCKER_CONTAINER = ${DOCKER_IMAGE}:${TAGGED_VERSION} 

.PHONY: build tag push test docker

build: 
	docker build -f Dockerfile --tag ${DOCKER_CONTAINER} .

tag: 
	docker tag ${DOCKER_CONTAINER} ${DOCKER_USERNAME}/${DOCKER_CONTAINER}

push: 
	docker push ${DOCKER_USERNAME}/${DOCKER_CONTAINER}

test: 
	go test ./...

docker: test build