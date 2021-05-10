NAME=Authentication_Server
VERSION=0.0.1

ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

.PHONY: docker
docker:
	@docker run -v "$(ROOT_DIR)/src":"/web" -p 8080:8080 --name "auth_server_dev" --rm ldesplanche/auth_server_dev

.PHONY: docker-dev-image
docker-dev-image:
	@docker build -t ldesplanche/auth_server_dev - < dev.Dockerfile

.PHONY: build
## build: Compile the packages.
build:
	@go build -o $(NAME)

.PHONY: run
## run: Build and Run in development mode.
run: build
	@./$(NAME) -e dev

.PHONY: run-prod
## run-prod: Build and Run in production mode.
run-prod: build
	@./$(NAME) -e prod

.PHONY: clean
## clean: Clean project and previous builds.
clean:
	@rm -f $(NAME)

.PHONY: deps
## deps: Download modules
deps:
	@go mod download

.PHONY: test
## test: Run tests with verbose mode
test:
	@go test -v ./tests/*

.PHONY: help
all: help
# help: show this help message
help: Makefile
	@echo
	@echo " Choose a command to run in "$(APP_NAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo