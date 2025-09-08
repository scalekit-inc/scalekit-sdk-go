.PHONY: generate install-tools

# Generate protobuf files
generate:
	buf generate --template buf.gen.yaml proto --verbose

# Install required tools
install-tools:
	go install github.com/bufbuild/buf/cmd/buf@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.6
	go install connectrpc.com/connect/cmd/protoc-gen-connect-go@latest