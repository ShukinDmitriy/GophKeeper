FROM golang:1.22

WORKDIR /app

RUN go install golang.org/x/tools/cmd/goimports@latest

CMD ["/app/docker/server/docker-entrypoint.sh"]
