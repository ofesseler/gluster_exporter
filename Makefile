GO = go

info:
	@echo "build: Go build"
	@echo "docker: build and run in docker container"

build: gotest
	$(GO) build -o gluster_exporter main.go

docker: gotest build
	docker build -t gluster-exporter-test .
	docker run --rm --privileged=true -p 9189:9189 -p 24007:24007 -p 24008:24008 -ti -v gluster-test:/data gluster-exporter-test

gotest:
	$(GO) test -v .
