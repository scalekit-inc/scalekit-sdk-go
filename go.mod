module github.com/scalekit-inc/scalekit-sdk-go

go 1.22.0

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.33.0-20240221180331-f05a6f4403ce.1
	connectrpc.com/connect v1.16.0
	github.com/go-jose/go-jose/v4 v4.0.1
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.19.1
	google.golang.org/genproto/googleapis/api v0.0.0-20240325203815-454cdb8f5daa
	google.golang.org/protobuf v1.33.0
)

require golang.org/x/crypto v0.21.0 // indirect
