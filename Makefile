GO ?= go

.PHONY: all
all: build

.PHONY: build
build:
	$(GO) install suse.com/caaspctl/cmd/...

.PHONY: staging
staging:
	$(GO) install -tags staging suse.com/caaspctl/cmd/...

.PHONY: release
release:
	$(GO) install -tags release suse.com/caaspctl/cmd/...
