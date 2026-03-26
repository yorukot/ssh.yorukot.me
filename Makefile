run:
	go run cmd/main.go

dev: 
	air go run cmd/main.go

lint:
	go fmt ./...
	go vet ./...
	golangci-run run