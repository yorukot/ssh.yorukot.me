FROM golang:1.25-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/ssh.yorukot.me ./cmd/main.go

FROM alpine:3.22

RUN apk add --no-cache openssh-keygen

WORKDIR /app

COPY --from=builder /out/ssh.yorukot.me /usr/local/bin/ssh.yorukot.me
COPY content ./content
COPY yorukot.me/src/content/blog ./content/markdown/blog

EXPOSE 23234

CMD ["sh", "-c", "mkdir -p .ssh && [ -f .ssh/id_ed25519 ] || ssh-keygen -t ed25519 -f .ssh/id_ed25519 -N '' && exec ssh.yorukot.me"]
