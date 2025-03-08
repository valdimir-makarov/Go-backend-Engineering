# Variables
PROTO_DIR := protobuf
GO_OUT := services/common/genproto/orders

# Generate Go code from .proto files
generate:
	# Ensure the output directory exists
	mkdir -p $(GO_OUT)
	protoc \
		--proto_path=$(PROTO_DIR) \
		--go_out=$(GO_OUT) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(GO_OUT) \
		--go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/orders.proto
	@echo "Code generation complete."
