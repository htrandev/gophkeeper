ifndef ROOT
	ROOT = $(shell pwd)
endif

# Если аргумент proto_out не указан, то он принимает значение по умолчанию.
ifndef proto_out
override proto_out = .
endif

# Если аргумент proto_pkg не указан, то он принимает значение по умолчанию.
ifndef proto_pkg
override proto_pkg = ./proto;proto
endif

.PHONY: proto
proto: | $(PROTOC) $(PROTOC_GEN_GO) $(PROTOC_GEN_GO_GRPC) ## Генерирует api/*.proto
	@protoc --proto_path=$(ROOT) --go_out=$(proto_out) --go-grpc_out=$(proto_out) \
		$(shell ls $(ROOT)/api/ | sed "s%.*%--go_opt=M'api/&=$(proto_pkg)' --go-grpc_opt=M'api/&=$(proto_pkg)'%") \
		$(ROOT)/api/*.proto

.PHONY: migration.sql
migration.sql:
	@goose -dir ./internal/storage/postgres/migrations create $(name) sql 

.PHONY: build.server
build.server:
	@go build -o cmd/server/server cmd/server/*.go


ifndef db_dsn
override db_dsn = postgres://user:password@localhost:25434/gophkeeper
endif

.PHONY: run.server
run.server:
	@DATABASE_DSN=$(db_dsn) go run ./cmd/server/...
