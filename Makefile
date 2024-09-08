# Created by Anh Cao on 27.08.2024.

APPLICATION_NAME ?= electric-push-notifications
TAG_VERSION ?= latest

docker: 
	go test ./... && docker build --tag ${APPLICATION_NAME}/${TAG_VERSION} .
