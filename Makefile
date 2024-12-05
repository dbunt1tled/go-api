run_api:
	@go run cmd/api/main.go
build_api:
	@go build -o bin/api cmd/api/main.go
build_opt_api:
	@go build -ldflags "-s -w"  -o bin/api cmd/api/main.go
gen_proto:
	@protoc --proto_path=proto proto/*.proto  --go-grpc_out=../ --go_out=../
gen_clean:
	@rm -rf ./proto/*.pb.go