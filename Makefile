GIT_COMMIT=$(shell git rev-parse HEAD)
GIT_DIRTY=$(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)
BUILD_TIME=$(shell date -u '+%Y-%m-%dT%I:%M:%H%z')

LD_FLAGS="-X main.GitCommit=${GIT_COMMIT}${GIT_DIRTY} -X main.BuildTime=${BUILD_TIME}"

.PHONY: build

build:
	go build -ldflags ${LD_FLAGS}

default: build
