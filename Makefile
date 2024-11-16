run_api:
	@go run cmd/api/main.go
build_api:
	@go build -o bin/api cmd/api/main.go
