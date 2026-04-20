run:
	go run cmd/main.go

dev:
	air go run cmd/main.go

build:
	go build -o bin/ssh.yorukot.me cmd/main.go 

update-blog:
	git -C yorukot.me fetch origin main
	git -C yorukot.me reset --hard origin/main
	$(MAKE) generate-blog-image-manifest
	git add yorukot.me
	git add content/blog_image_manifest.json
	git commit -m "Update blog content"

generate-blog-image-manifest:
	pnpm --dir yorukot.me install --frozen-lockfile
	pnpm --dir yorukot.me build
	go run ./cmd/blogimagemanifest

lint:
	go fmt ./...
	go vet ./...
	golangci-run run
