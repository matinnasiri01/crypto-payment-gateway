all: build test
.PHONY: all build test watch migrateup migratedown
include .env


migrateup:
	@migrate -path ./migrations/ -database "${DATABASE_URL}?sslmode=disable" up

migratedown:
	@migrate -path ./migrations/ -database "${DATABASE_URL}?sslmode=disable" down


build:
	@go build -o main.exe cmd/api/main.go

test:
	@go test ./... -v


watch:
	@powershell -ExecutionPolicy Bypass -Command "if (Get-Command air -ErrorAction SilentlyContinue) { \
		air; \
		Write-Output 'Watching...'; \
	} else { \
		Write-Output 'Installing air...'; \
		go install github.com/air-verse/air@latest; \
		air; \
		Write-Output 'Watching...'; \
	}"

