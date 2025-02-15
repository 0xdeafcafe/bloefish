# Build the Go code
build-go:
	rm -rf ./bin
	mkdir -p ./bin
	GOOS=linux go build -o ./bin/bloefish ./cmd/bloefish/...

build-js:
	yarn build

build-docker:
	docker build -t go_service -f go_service.Dockerfile .

install:
	yarn
	go mod download

build: build-js build-go build-docker

start:
	docker-compose up --build -d

# Aliases for build&start, but only does what is necessary
services: build-go docker-build start
web: build-js docker-build start
infra: docker-build start
