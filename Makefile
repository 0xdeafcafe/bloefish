# Build the Go code
build-go:
	rm -rf ./bin
	mkdir -p ./bin
	GOOS=linux go build -o ./bin/bloefish ./cmd/bloefish/...

# Build the JS code
build-js:
	yarn build

# Build the Docker image
docker-build:
	docker build -t go_service -f go_service.Dockerfile .

# Start Docker Compose
docker-up:
	docker-compose up --build -d

# Combined command to build Go code, build Docker image, and start Docker Compose
all: build-js build-go docker-build docker-up

services: build-go docker-build docker-up
web: build-js docker-build docker-up
infra: docker-build docker-up
