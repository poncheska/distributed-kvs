linters:
	golangci-lint run ./... -v

test:
	go test --race ./...

swagger:
	swag init -g internal/server/http/server.go --parseInternal --pd

up:
	docker-compose up --force-recreate --build app nginx


.PHONY: .vendor-rm
.vendor-rm:
	rm -rf vendor-proto

# Устанавливаем proto описания google/protobuf
vendor-proto/google/protobuf:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
		https://github.com/protocolbuffers/protobuf vendor-proto/protobuf &&\
	cd vendor-proto/protobuf &&\
	git sparse-checkout set --no-cone src/google/protobuf &&\
	git checkout
	mkdir -p vendor-proto/google
	mv vendor-proto/protobuf/src/google/protobuf vendor-proto/google
	rm -rf vendor-proto/protobuf


# Вендоринг внешних proto файлов
.vendor-proto: .vendor-rm  vendor-proto/google/protobuf

# Используем bin в текущей директории для установки плагинов protoc
LOCAL_BIN:=$(CURDIR)/bin

.PHONY: .bin-deps
.bin-deps:
	$(info Installing binary dependencies...)

	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1 && \
    GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0

STORE_PROTO_PATH:="api/store"

PHONY: .protoc-generate
.protoc-generate: .bin-deps .vendor-proto
	protoc \
	-I ${STORE_PROTO_PATH} \
	-I vendor-proto \
	--plugin=protoc-gen-go=$(LOCAL_BIN)/protoc-gen-go \
	--go_out pkg/${STORE_PROTO_PATH} \
	--go_opt paths=source_relative \
	--plugin=protoc-gen-go-grpc=$(LOCAL_BIN)/protoc-gen-go-grpc \
	--go-grpc_out pkg/${STORE_PROTO_PATH} \
	--go-grpc_opt paths=source_relative \
	api/store/store.proto
	go mod tidy