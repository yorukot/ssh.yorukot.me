run:
	go run cmd/main.go

dev:
	air go run cmd/main.go

build:
	go build -o bin/ssh.yorukot.me cmd/main.go 

lint:
	go fmt ./...
	go vet ./...
	golangci-run run
