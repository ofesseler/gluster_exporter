GO = go
pkgs = $(shell $(GO) list ./... | grep -v /vendor/)
PROMU = promu
BIN_DIR ?= $(shell pwd)

info:
	@echo "build: Go build"
	@echo "docker: build and run in docker container"
	@echo "gotest: run go tests and reformats"

build: gotest
	$(PROMU) build
	#$(GO) build -o gluster_exporter

docker: gotest build
	docker build -t gluster-exporter-test .
	docker run --rm --privileged=true -p 9189:9189 -p 24007:24007 -p 24008:24008 -i -v gluster-test:/data gluster-exporter-test

gotest: fmt
	$(GO) test -v $(pkgs)

fmt:
	$(GO) fmt

promu:
	@GOOS=$(shell uname -s | tr A-Z a-z) \
		GOARCH=$(subst x86_64,amd64,$(patsubst i%86,386,$(shell uname -m))) \
		$(GO) get -u github.com/prometheus/promu

promu-build: gotest promu
	$(PROMU) build

tarball: build promu
	@$(PROMU) tarball $(BIN_DIR)
