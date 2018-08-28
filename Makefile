SVC_NAME=rabbit_consumer_groups
# VERSION=`git describe --abbrev=0 --tags`
VERSION=latest

DOCKER_REGISTRY=kskitek
DOCKER_IMAGE=$(DOCKER_REGISTRY)/$(SVC_NAME):$(VERSION)

all: help

## build: builds project binary
build:
	GO111MODULE=on go1.11rc2 build ./...

build-linux:
	env GOOS=linux GO11MODULE=on go build -o $(SVC_NAME)_linux

## run: runs project locally
run: build
	env $$(cat config/service.env) ./$(SVC_NAME)

## deps: tidies and downloads all dependencies
deps:
	GO111MODULE=on go1.11rc2 mod tidy
	GO111MODULE=on go1.11rc2 mod download

update-deps:
	GO111MODULE=on go1.11rc2 get -u

docker-build: build-linux
	docker build -t $(DOCKER_IMAGE) .

docker-push: docker-build
	docker push $(DOCKER_IMAGE)

## docker: builds and pushes new image
docker: docker-build docker-push

## docker-run: runs project from docker-compose
docker-run: docker-build
	docker-compose up


help: Makefile
	@echo " Choose a command run in \033[32m"$(SVC_NAME)"\033[0m:"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'