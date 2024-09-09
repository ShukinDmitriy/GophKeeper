#!/bin/sh

printf "Install packages...\n\n"
go mod tidy

printf "Apply migrations...\n\n"
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
/usr/bin/make -f Makefile migrate-up

printf "Run tests...\n\n"
cd /GophKeeper
/usr/bin/make test-cover
