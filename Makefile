include .env

run_api:
	@go run cmd/api/api.go
build_api:
	@go build -ldflags "-s -w"  -o bin cmd/api/api.go
run_email_consumer:
	@go run cmd/email/email.go
build_email_consumer:
	@go build -ldflags "-s -w"  -o bin cmd/email/email.go
migrate:
	@go run cmd/migrator/migrator.go
install_govulncheck:
	@go install golang.org/x/vuln/cmd/govulncheck@latest
check_vulnerabilities:
	@govulncheck ./...
run_centrifugo:
	@go run cmd/centrifugo/centrifugo_server.go
build_centrifugo:
	@go build -ldflags "-s -w"  -o bin cmd/centrifugo/centrifugo_server.go
gen_proto:
	@protoc --proto_path=proto proto/*.proto  --go-grpc_out=./internal/grpc/proxyproto --go_out=./internal/grpc/proxyproto
gen_clean:
	@rm -rf ./internal/grpc/proxyproto/*.pb.go

.PHONY: run_api build_api,run_email_consumer,build_email_consumer,migrate gen_proto gen_clean, run_centrifugo,build_centrifugo
