#!/bin/sh

printf "Install packages...\n\n"
go mod tidy

printf "Run tests...\n\n"
cd /app \
&& make test-cover
#   -accrual-database-uri="postgres://postgres:postgres@postgres/praktikum?sslmode=disable"
