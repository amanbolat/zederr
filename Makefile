export BIN := $(PWD)/.bin
export PATH := $(BIN):$(PATH)

gen-proto:
	# install protoc
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0
	protoc --proto_path=. --go_out=. --go-grpc_out=. \
	 --go_opt=module=github.com/amanbolat/zederr \
	 --go-grpc_opt=module=github.com/amanbolat/zederr  \
	 pkg/proto/v1/*.proto

.PHONY: bin
bin:
	mkdir -p .bin

.PHONY: bin.go-enum
bin.go-enum: bin
	cd tools/deps && go mod tidy && GOBIN=$(BIN) go install -modfile go.mod github.com/abice/go-enum

.PHONY: gen.enums
gen.enums: bin.go-enum
	go-enum -file internal/codegen/core/argument_type.go --marshal --sql --nocase
