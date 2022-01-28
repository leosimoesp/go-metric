include .env
export

.PHONY: start

run:
	go run cmd/main.go

create-mocks:
	go generate ./...

.PHONY: debug
debug:
	dlv debug --headless --listen=:2345 --log --api-version=2

