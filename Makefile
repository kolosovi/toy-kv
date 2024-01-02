CWD=$(shell pwd)

.PHONY: bin-deps
bin-deps:
	GOBIN="${CWD}/bin" go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.32.0

.PHONY: generate
generate:
	PATH="${PATH}:${CWD}/bin" protoc -I=internal/proto/wal --go_out=. internal/proto/wal/wal.proto