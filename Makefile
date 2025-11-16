APP_NAME=reviewer-app

.PHONY: docker-up docker-down docker-restart build run

docker-up:
	docker compose up -d --build

docker-down:
	docker compose down

docker-restart:
	docker compose down
	docker compose up -d --build

run:
	go run ./cmd/app/main.go

build:
	go build -o $(APP_NAME) ./cmd/app/main.go
