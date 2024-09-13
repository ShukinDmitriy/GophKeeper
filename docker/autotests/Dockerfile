FROM golang:1.22

WORKDIR /GophKeeper

RUN go install golang.org/x/tools/cmd/goimports@latest

ENTRYPOINT ["/GophKeeper/docker/autotests/docker-entrypoint.sh"]