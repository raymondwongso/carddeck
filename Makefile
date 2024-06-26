GO_VERSION 						= $(shell awk '$$1 ~ /^go$$/ {print $$2}' go.mod)
GOIMPORTS_VERSION 		= v0.20.0
GOWORKING_DIR 				= src/app
GOLANGCILINT_VERSION 	= v1.57.2
GO_FILES							= $(shell go list ./... | grep -v /docs | grep -v /config)

.PHONY: format
format:
	docker run --rm -v ./:/go/$(GOWORKING_DIR) golang:$(GO_VERSION) sh -c "go install golang.org/x/tools/cmd/goimports@$(GOIMPORTS_VERSION); cd $(GOWORKING_DIR); goimports -w .; gofmt -s -w ."

.PHONY: lint
lint:
	docker run --rm -v ./:/app -v ~/.cache/golangci-lint/$(GOLANGCILINT_VERSION):/root/.cache -w /app golangci/golangci-lint:$(GOLANGCILINT_VERSION) golangci-lint run -v

.PHONY: mockgen
mockgen:
	bin/generate-mock.sh

.PHONY: test
test:
	bin/test.sh

.PHONY: swagger
swagger:
	bin/swagger.sh

.PHONY: build
build:
	go build .

.PHONY: run
run: build
	./carddeck server
