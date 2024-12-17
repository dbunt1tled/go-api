include .env

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

#  MIGRATION_NAME=create_table_users make migration_sql
migration_sql:
	@GOOSE_DRIVER="${GOOSE_DRIVER}" GOOSE_DBSTRING="${GOOSE_DBSTRING}" goose -dir ./internal/database/migrations create $(MIGRATION_NAME) sql
migration_go:
	@GOOSE_DRIVER="${GOOSE_DRIVER}" GOOSE_DBSTRING="${GOOSE_DBSTRING}" goose -dir ./internal/database/migrations create $(MIGRATION_NAME) go
migrate_up:
	@GOOSE_DRIVER="${GOOSE_DRIVER}" GOOSE_DBSTRING="${GOOSE_DBSTRING}" goose -dir ./internal/database/migrations up
migrate_down:
	@GOOSE_DRIVER="${GOOSE_DRIVER}" GOOSE_DBSTRING="${GOOSE_DBSTRING}" goose -dir ./internal/database/migrations down
migrate_status:
	@GOOSE_DRIVER="${GOOSE_DRIVER}" GOOSE_DBSTRING="${GOOSE_DBSTRING}" goose -dir ./internal/database/migrations status

.PHONY: run_api build_api build_opt_api gen_proto gen_clean migration_sql migration_go migrate_up migrate_down migrate_status