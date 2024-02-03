gen-proto:
	# install protoc
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0
	protoc --proto_path=. --go_out=. --go-grpc_out=. \
	 --go_opt=module=github.com/amanbolat/zederr \
	 --go-grpc_opt=module=github.com/amanbolat/zederr  \
	 proto/v1/*.proto

