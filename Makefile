linters:
	golangci-lint run ./... -v

test:
	go test --race ./...

swagger:
	swag init -g internal/server/http/server.go --parseInternal --pd

up:
	docker-compose up --force-recreate --build app nginx