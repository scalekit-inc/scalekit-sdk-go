module github.com/scalekit-inc/scalekit-sdk-go

go 1.22.0

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.34.2-20240717164558-a6c49f84cc0f.2
	connectrpc.com/connect v1.16.2
	github.com/go-jose/go-jose/v4 v4.0.5
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.23.0
	google.golang.org/genproto/googleapis/api v0.0.0-20241021214115-324edc3d5d38
	google.golang.org/protobuf v1.35.2
)

require golang.org/x/crypto v0.32.0 // indirect
