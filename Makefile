SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

BINARY=staticli
BUILD_TIME=`date +%FT%T%z`

DEFAULT_SYSTEM_BINARY := $(BINARY).darwin.amd64
BINTRAY_API_KEY=$(shell cat api_key)
GITHUB_API_KEY=$(shell cat github_api)
VERSION=$(shell cat VERSION)
BUILD_TIME=$(shell date +%FT%T%z)
BUILD_COMMIT=$(shell git rev-parse HEAD)
BUILD_REPO=$(shell git remote get-url origin)

UNAME_S := $(shell uname -s)
DEFAULT_SHASUM_UTIL=shasum
ifeq ($(UNAME_S),Linux)
	DEFAULT_SHASUM_UTIL=sha1sum
	DEFAULT_SYSTEM_BINARY := $(BINARY).linux.amd64
endif

ifndef TRAVIS
	DOCKER_RUN_COMMAND=docker run --rm -v $(shell pwd)/../../../:/go/src/ -w /go/src/github.com/staticli/staticli
endif
ifdef TRAVIS
	DOCKER_RUN_COMMAND=docker run --rm -v $(shell pwd)/:/go/src/github.com/staticli/staticli -w /go/src/github.com/staticli/staticli
endif

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS=-ldflags \"-X github.com/staticli/staticli/lib.Version=${VERSION} -X github.com/staticli/staticli/lib.BuildTime=${BUILD_TIME} -X github.com/staticli/staticli/lib.BuildCommit=${BUILD_COMMIT} -X github.com/staticli/staticli/lib.BuildRepo=${BUILD_REPO}\"

.DEFAULT_GOAL: $(BINARY)
$(BINARY): $(BINARY).darwin.amd64 $(BINARY).linux.amd64 $(BINARY).linux.arm
	cp $(DEFAULT_SYSTEM_BINARY) $@

$(BINARY).darwin.amd64: $(SOURCES)
	${DOCKER_RUN_COMMAND} -e GOOS=darwin -e GOARCH=amd64 staticli/godep /bin/bash -c "dep ensure && go build ${LDFLAGS} -o $@"
	${DEFAULT_SHASUM_UTIL} $@ > $@.sha

$(BINARY).linux.amd64: $(SOURCES)
	${DOCKER_RUN_COMMAND} -e GOOS=linux -e GOARCH=amd64 staticli/godep /bin/bash -c "dep ensure && go build ${LDFLAGS} -o $@"
	${DEFAULT_SHASUM_UTIL} $@ > $@.sha

$(BINARY).linux.arm: $(SOURCES)
	${DOCKER_RUN_COMMAND} -e GOOS=linux -e GOARCH=arm staticli/godep /bin/bash -c "dep ensure && go build ${LDFLAGS} -o $@"
	${DEFAULT_SHASUM_UTIL} $@ > $@.sha


.PHONY: clean
clean:
	rm -f -- ${BINARY}
	rm -f -- ${BINARY}.darwin.amd64 ${BINARY}.darwin.amd64.sha
	rm -f -- ${BINARY}.linux.amd64  ${BINARY}.linux.amd64.sha
	rm -f -- ${BINARY}.linux.arm    ${BINARY}.linux.arm.sha

.PHONY: install
install:
	cp $(DEFAULT_SYSTEM_BINARY) ~/bin/staticli

.PHONY: release
release:
	./staticli github-release "${VERSION}" staticli.darwin.amd64 staticli.linux.amd64 staticli.linux.arm -- --github-access-token ${GITHUB_API_KEY} --github-repository staticli/staticli

# Really simple "does it at least run?" tests for now
# Proper tests coming at some point
.PHONY: test
test: test-unit test-integration test-binary

.PHONY: test-unit
test-unit:
	echo "Coming soon"

.PHONY: test-integration
test-integration:
	go run main.go -d version
	echo "More coming at some point"

.PHONY: test-binary
test-binary:
	./$(DEFAULT_SYSTEM_BINARY) version