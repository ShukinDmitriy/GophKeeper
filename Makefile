include server.env

migrate-create:
	@(printf "Enter migrate name: "; read arg; migrate create -ext sql -dir db/migrations -seq $$arg);

migrate-up:
	migrate -database ${DATABASE_URI} -path ./db/migrations up

migrate-down:
	migrate -database ${DATABASE_URI} -path ./db/migrations down 1

migrate-down-all:
	migrate -database ${DATABASE_URI} -path ./db/migrations down --all

build-mocks:
	@go get github.com/vektra/mockery/v2@v2.43.2
	@~/go/bin/mockery

server:
	docker compose -f server.docker-compose.yml start

test:
	docker compose -f test.docker-compose.yml up -d
	docker compose -f test.docker-compose.yml logs -f golang
	docker compose -f test.docker-compose.yml down

test-cover:
	go test -v -coverprofile=coverage.out ./... && go tool cover -html=coverage.out -o coverage.html

test-cover-cli:
	./test-cover.sh

static-check:
	go vet ./cmd/... ./internal/...

format:
	@go get mvdan.cc/gofumpt@latest
	@~/go/bin/gofumpt -l -w .

generate-swagger:
	~/go/bin/swag init -g ./cmd/server/main.go -o ./cmd/server/docs