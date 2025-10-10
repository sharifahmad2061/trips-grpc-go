PROTOC_GEN_GO_VERSION = v1.36.10
PROTOC_GEN_GO_GRPC_VERSION = v1.5.1

.PHONY: tools
tools:
	brew install protobuf
	go install google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO_VERSION)
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GEN_GO_GRPC_VERSION)

PROTO_DIR = api/proto
OUT_DIR = api/gen/go

.PHONY: gen_proto
gen_proto:
	protoc -I $(PROTO_DIR) --go_out=$(OUT_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(OUT_DIR) --go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/*.proto

.PHONY: build
build: go build