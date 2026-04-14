BINARY    := clash-forge
MODULE    := github.com/wuruimiao/clash-forge
CMD       := ./cmd/clash-forge
DIST      := dist

VERSION   ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT    := $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
BUILD_AT  := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')

LDFLAGS   := -s -w \
             -X 'main.version=$(VERSION)'
GOFLAGS   := -trimpath

# 默认目标平台
PLATFORMS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64

.PHONY: build build-all clean test lint help

## build: 编译当前平台二进制（输出到 dist/）
build:
	@mkdir -p $(DIST)
	CGO_ENABLED=0 go build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(DIST)/$(BINARY) $(CMD)
	@echo "Built: $(DIST)/$(BINARY)"

## build-all: 交叉编译所有目标平台
build-all:
	@mkdir -p $(DIST)
	@for platform in $(PLATFORMS); do \
		os=$${platform%/*}; \
		arch=$${platform#*/}; \
		output=$(DIST)/$(BINARY)-$${os}-$${arch}; \
		if [ "$${os}" = "windows" ]; then output=$${output}.exe; fi; \
		echo "Building $${os}/$${arch} ..."; \
		CGO_ENABLED=0 GOOS=$${os} GOARCH=$${arch} \
			go build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $${output} $(CMD) || exit 1; \
	done
	@echo "All builds complete:"
	@ls -lh $(DIST)/

## test: 运行所有测试
test:
	go test ./... -count=1

## lint: 运行 go vet
lint:
	go vet ./...

## clean: 清理构建产物
clean:
	rm -rf $(DIST)

## help: 显示帮助
help:
	@grep -E '^## ' Makefile | sed 's/## //' | column -t -s ':'
