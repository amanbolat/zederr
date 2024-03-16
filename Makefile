export BIN := $(PWD)/.bin
export PROTOC_BIN := $(PWD)/.bin/protocd/bin
export PATH := $(BIN):$(PROTOC_BIN):$(PATH)
PROTOC_VER := 25.3
PROTOC_CHECKSUM := d0fcd6d3b3ef6f22f1c47cc30a80c06727e1eccdddcaf0f4a3be47c070ffd3fe

gen-proto:
	# install protoc
	PROTOC_VER=${PROTOC_VER} PROTOC_CHECKSUM=${PROTOC_CHECKSUM} ./scripts/install-protoc.sh
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	protoc --proto_path=. --go_out=. --go-grpc_out=. \
	 --go_opt=module=github.com/amanbolat/zederr \
	 zeproto/v1/*.proto

.PHONY: bin
bin:
	mkdir -p .bin

.PHONY: bin.go-enum
bin.go-enum: bin
	cd tools/deps && go mod tidy && GOBIN=$(BIN) go install -modfile go.mod github.com/abice/go-enum

.PHONY: bin.golangci-lint
bin.golangci-lint: bin
	cd tools/deps && go mod tidy && GOBIN=$(BIN) go install -modfile go.mod github.com/golangci/golangci-lint/cmd/golangci-lint

.PHONY: lint
lint: bin.golangci-lint
	$(BIN)/golangci-lint run

.PHONY: gen.enums
gen.enums: bin.go-enum
	$(BIN)/go-enum -file internal/codegen/core/argument_type.go --marshal --sql --nocase
