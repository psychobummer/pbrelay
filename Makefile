proto:
	protoc --go_out=plugins=grpc:. ./rpc/relay/relay.proto
	protoc --go_out=plugins=grpc:. ./rpc/keystore/keystore.proto

test:
	go test -v ./...

race:
	go test -race -v ./...
