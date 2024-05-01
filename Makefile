proto:
	protoc --go_out=./proto --go-grpc_out=./proto ./proto/types.proto

.PHONY: proto