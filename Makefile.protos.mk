################################################################################
## Protobuf related targets
################################################################################

PROTOTOOL ?= $(GORUN) github.com/uber/prototool/cmd/prototool

.PHONY: protos
protos: plugins protos/compile protos/generate

.PHONY: protos/compile
protos/compile:
	$(PROTOTOOL) compile

.PHONY: protos/generate
protos/generate:
	env PATH="bin:$$PATH" $(PROTOTOOL) generate

.PHONY: plugins
plugins: plugins/go plugins/go-grpc plugins/grpc-gateway

.PHONY: plugins/go
plugins/go:
	GOBIN="$$PWD/bin" go install google.golang.org/protobuf/cmd/protoc-gen-go

.PHONY: plugins/go-grpc
plugins/go-grpc:
	GOBIN="$$PWD/bin" go install google.golang.org/grpc/cmd/protoc-gen-go-grpc

.PHONY: plugins/grpc-gateway
plugins/grpc-gateway:
	GOBIN="$$PWD/bin" go install \
		github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway

