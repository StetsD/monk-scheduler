protocomplile:
	protoc ./internal/api/event.proto --go_out=plugins=grpc:.

build:
	go build -o monkapp main.go

setup:
	go run main.go setup

start:
	$(MAKE) --no-print-directory lint
	-go run main.go start

test:
	go test -v ./tests/...

migrate:
	flyway migrate

run:
	$(MAKE) --no-print-directory setup
	$(MAKE) --no-print-directory migrate
	-$(MAKE) --no-print-directory start

lint:
	golangci-lint run