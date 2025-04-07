.PHONY: lint
lint:
	@golangci-lint -v run --fix ./...x

.PHONY: imports
imports:
	@goimports-reviser -project-name github.com/AFK068/bot -file-path ./... -separate-named

.PHONY: generate_openapi
generate_openapi:
	@mkdir -p internal/api/openapi/compressor/v1
	@oapi-codegen -package v1 \
		-generate server,types \
		api/openapi/v1/compressor.yaml > internal/api/openapi/compressor/v1/compressor-api.gen.go