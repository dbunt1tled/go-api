include .env

run_api:
	@go run cmd/api/api.go
build_api:
	@go build -o bin cmd/api/api.go
build_opt_api:
	@go build -ldflags "-s -w"  -o bin cmd/api/api.go
install_govulncheck:
	@go install golang.org/x/vuln/cmd/govulncheck@latest
check_vulnerabilities:
	@govulncheck ./...

# MIGRATION_NAME=create_table_user make migration_sql
migration_sql:
	@GOOSE_DRIVER="${GOOSE_DRIVER}" GOOSE_DBSTRING="${GOOSE_DBSTRING}" GOOSE_MIGRATION_DIR="${GOOSE_MIGRATION_DIR}" goose create $(MIGRATION_NAME) sql
migration_go:
	@GOOSE_DRIVER="${GOOSE_DRIVER}" GOOSE_DBSTRING="${GOOSE_DBSTRING}" GOOSE_MIGRATION_DIR="${GOOSE_MIGRATION_DIR}" goose create $(MIGRATION_NAME) go
migrate_up:
	@GOOSE_DRIVER="${GOOSE_DRIVER}" GOOSE_DBSTRING="${GOOSE_DBSTRING}" GOOSE_MIGRATION_DIR="${GOOSE_MIGRATION_DIR}" goose  up
migrate_down:
	@GOOSE_DRIVER="${GOOSE_DRIVER}" GOOSE_DBSTRING="${GOOSE_DBSTRING}" GOOSE_MIGRATION_DIR="${GOOSE_MIGRATION_DIR}" goose down
migrate_status:
	@GOOSE_DRIVER="${GOOSE_DRIVER}" GOOSE_DBSTRING="${GOOSE_DBSTRING}" GOOSE_MIGRATION_DIR="${GOOSE_MIGRATION_DIR}" goose status

.PHONY: run_api build_api build_opt_api gen_proto gen_clean migration_sql migration_go migrate_up migrate_down migrate_status
