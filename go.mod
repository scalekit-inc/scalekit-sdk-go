module github.com/scalekit-inc/scalekit-sdk-go

go 1.22.0

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.34.2-20240508200655-46a4cf4ba109.2
	connectrpc.com/connect v1.16.2
	github.com/go-jose/go-jose/v4 v4.0.4
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.21.0
	google.golang.org/genproto/googleapis/api v0.0.0-20240723171418-e6d459c13d2a
	google.golang.org/protobuf v1.34.2
)

require golang.org/x/crypto v0.25.0 // indirect
