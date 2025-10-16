.PHONY: proto

proto:
	@MOD="github.com/MystenLabs/sui-apis/proto"; \
	GO_M_OPTS=$$(for f in proto/sui/rpc/v2/*.proto; do \
		imp=$${f#proto/}; \
		printf -- "--go_opt=M%s=$${MOD}/sui/rpc/v2 " "$$imp"; \
	done); \
	GRPC_M_OPTS=$${GO_M_OPTS//--go_opt=/--go-grpc_opt=}; \
	eval protoc -I ./proto \
		--go_out=paths=source_relative:. $$GO_M_OPTS \
		--go-grpc_out=paths=source_relative:. $$GRPC_M_OPTS \
		proto/sui/rpc/v2/*.proto
