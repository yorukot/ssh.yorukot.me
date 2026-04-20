run:
	go run cmd/main.go

dev:
	air go run cmd/main.go

build:
	go build -o bin/ssh.yorukot.me cmd/main.go 

update-blog:
	git -C yorukot.me pull origin main
	git add yorukot.me
	git commit -m "Update blog content"

lint:
	go fmt ./...
	go vet ./...
	golangci-run run
