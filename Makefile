CMD := docker-compose
SRCS := main.go $(wildcard ./**/*.go)

all:
	@echo "make up    :: run docker containers"
	@echo "make build :: build docker containers"
	@echo "make run   :: run the application"
	@echo "make down  :: stop docker containers"
	@echo "make fmt   :: format source files;"
	@echo ""
	@echo "managed source files;"
	@echo $(SRCS)

.PHONY: up build down format clean
up:
	@$(CMD) up -d

build:
	@$(CMD) build 

run:
	@$(CMD) exec app main

down:
	@$(CMD) down

fmt:
	go fmt ./...

clean:
	@$(CMD) down --rmi all --volumes --remove-orphans
