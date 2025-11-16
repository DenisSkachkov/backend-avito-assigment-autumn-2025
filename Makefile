# APP_NAME=reviewer-app
# APP_CMD=./cmd/main.go

# DB_URL=postgres://postgres:postgres@localhost:5432/reviewer?sslmode=disable

# MIGRATE=migrate -path migrations -database "$(DB_URL)"

# .PHONY: run build migrate-up migrate-down docker-build docker-up docker-down

# run:
# 	go run $(APP_CMD)

# build:
# 	go build -o $(APP_NAME) $(APP_CMD)

# migrate-up:
# 	$(MIGRATE) up

# migrate-down:
# 	$(MIGRATE) down

# docker-build:
# 	docker build -t $(APP_NAME) .

# docker-up:
# 	docker compose up -d

# docker-down:
# 	docker compose down


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
