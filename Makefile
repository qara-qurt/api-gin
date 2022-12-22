.PHONY: build
build:
	go build -v ./cmd/main.go

.PHONY: run
run:
	go run ./cmd/main.go	

# createdb:
# 	docker exec -it todo-db createdb --username=root --owner=root todo
docker-app:
	docker run -d --name=go-app -p 8080:8080 --rm go-app
docker-postgres: 
	docker run --name=todo-db -e POSTGRES_USER='postgres' -e POSTGRES_PASSWORD='root' -p 5436:5432 -d --rm postgres:latest
migrate-up:
	migrate -path ./schema -database 'postgres://postgres:root@localhost:5436/postgres?sslmode=disable' up

migrate-down:
	migrate -path ./schema -database 'postgres://postgres:root@localhost:5436/postgres?sslmode=disable' down

swag:
	swag init -g cmd/main.go
.DEFAULT_GOAL := run
