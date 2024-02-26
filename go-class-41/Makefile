SOURCES := $(wildcard *.go cmd/*/*.go)

VERSION=$(shell git describe --tags --long --dirty 2>/dev/null)

## we must have tagged the repo at least once for VERSION to work
ifeq ($(VERSION),)
	VERSION = UNKNOWN
endif

# will make binary for program
sort: $(SOURCES)
	go build -ldflags "-X main.version=${VERSION}" -o $@ ./cmd/sort

.PHONY: lint
lint:
	golangci-lint run

.PHONY: committed
committed:
	@git diff --exit-code > /dev/null || (echo "** COMMIT YOUR CHANGES FIRST **"; exit 1)

# will make docker container
docker: $(SOURCES) build/Dockerfile
	docker build -t sort-anim:latest . -f build/Dockerfile --build-arg VERSION=$(VERSION)

# no account so can't be published
.PHONY: publish
publish: committed lint
	make docker
	docker tag  sort-anim:latest upreeti/sort-anim:$(VERSION)
	docker push upreeti/sort-anim:$(VERSION)