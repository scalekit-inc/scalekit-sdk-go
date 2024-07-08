module github.com/scalekit-inc/scalekit-sdk-go

go 1.22.0

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.34.2-20240508200655-46a4cf4ba109.2
	connectrpc.com/connect v1.16.2
	github.com/go-jose/go-jose/v4 v4.0.2
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.20.0
	google.golang.org/genproto/googleapis/api v0.0.0-20240701130421-f6361c86f094
	google.golang.org/protobuf v1.34.2
)

require golang.org/x/crypto v0.25.0 // indirect
