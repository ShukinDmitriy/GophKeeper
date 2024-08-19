include .env

build-mocks:
	@go get github.com/vektra/mockery/v2@v2.43.2
	@~/go/bin/mockery

server:
	docker compose -f server.docker-compose.yml start

test:
	docker compose -f test.docker-compose.yml start
	docker compose -f test.docker-compose.yml logs -f golang
	docker compose -f test.docker-compose.yml stop

test-cover:
	go test -v -coverprofile=coverage.out ./... && go tool cover -html=coverage.out -o coverage.html

test-cover-cli:
	./test-cover.sh

static-check:
	go vet ./cmd/... ./internal/...

format:
	@go get mvdan.cc/gofumpt@latest
	@~/go/bin/gofumpt -l -w .