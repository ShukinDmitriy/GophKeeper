#!/bin/sh

printf "Install packages...\n\n"
go mod tidy

printf "Run go...\n\n"
go run ./cmd/server
