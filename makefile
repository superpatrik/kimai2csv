GO := $(shell which go)
GOBUILD := $(GO) build

BINARY ?= kimai2csv

VERSION ?= $(shell git describe --tags --dirty --always 2>/dev/null || echo dev)
TARGETS := darwin/arm64 darwin/amd64 linux/amd64 linux/arm64 windows/amd64 windows/arm64

BUILD_DIR := build
DIST_DIR := $(BUILD_DIR)/dist

SRC := .
OUT ?= ./$(BUILD_DIR)/$(BINARY)

all: clean tidy release compile

clean:
	rm -rf build

tidy:
	$(GO) mod tidy

compile:
	echo "-> $(OUT)"
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 $(GOBUILD) -trimpath -ldflags="-s -w" -o $(OUT) $(SRC)

release: clean tidy
	mkdir -p "$(DIST_DIR)"
	for target in $(TARGETS); do \
		GOOS="$${target%/*}"; \
		GOARCH="$${target#*/}"; \
		NAME="$(BINARY)_$(VERSION)_$${GOOS}_$${GOARCH}"; \
		echo "Building $$NAME"; \
		$(MAKE) compile GOOS=$$GOOS GOARCH=$$GOARCH OUT=$(DIST_DIR)/$$NAME; \
		echo "Packing $$NAME.tar.gz"; \
		tar -C $(DIST_DIR) -czf $(DIST_DIR)/$$NAME.tar.gz $$NAME; \
	done

	echo "Creating checksums"
	( cd $(DIST_DIR) && sha256sum * > checksums.txt )

install:
	cp $(OUT) /usr/local/bin/$(BINARY)