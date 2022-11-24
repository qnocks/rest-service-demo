.SILENT:

build:
	go build -o ./.bin/app cmd/app/main.go

run: build
	./.bin/app

test:
	go test ./cmd/... ./internal/... -race -coverprofile=cover.out ./...
	make test.coverage

test.coverage:
	go tool cover -func=cover.out | grep "total"

lint:
	golangci-lint run

gen:
	mockgen -source=internal/service/service.go -destination=internal/service/mocks/mock.go
	mockgen -source=internal/storage/storage.go -destination=internal/storage/mocks/mock.go
