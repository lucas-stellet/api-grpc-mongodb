gen:
	protoc --proto_path=proto proto/*.proto --go_out=pkg/ --go-grpc_out=require_unimplemented_servers=false:pkg/

dev:
	air main.go

start:
	go run main.go

gobuild:
	GOOS=linux go build -o build/api