NAME	 := matterpoll-emoji
VERSION  := v0.0.2
REVISION := $(shell git rev-parse --short HEAD)

LDFLAGS := -ldflags="-s -w -X \"main.Version=$(VERSION)\" -X \"main.Revision=$(REVISION)\" -extldflags \"-static\""
DIST_DIRS := find * -type d -exec

.PHONY: glide deps clean check-style test coverage cross-build dist

all: dist test

glide:
ifeq ($(shell command -v glide 2> /dev/null),)
	curl https://glide.sh/get | sh
endif

deps: glide
	glide install

cross-build: deps
	for os in darwin linux windows; do \
		GOOS=$$os GOARCH=386 CGO_ENABLED=0 go build -a -tags netgo -installsuffix netgo $(LDFLAGS) -o dist/$$os-i686/$(NAME); \
	done 
	for os in darwin linux windows; do \
		GOOS=$$os GOARCH=amd64 CGO_ENABLED=0 go build -a -tags netgo -installsuffix netgo $(LDFLAGS) -o dist/$$os-x86_64/$(NAME); \
	done 

clean:
	rm -rf bin/*
	rm -rf vendor/*
	rm -rf dist/*

test:
	go test -v ./poll/

coverage:
	go test -coverprofile=coverage.txt -covermode=atomic ./poll/
	go tool cover -html=coverage.txt

check-style:
	@echo Running gofmt
	$(eval GOFILES_NOVENDOR = $(shell find . -type f -name '*.go' -not -path './vendor/*'))
	$(eval GOFMT_OUTPUT := $(shell gofmt -l -s $(GOFILES_NOVENDOR) 2>&1))
	@if [ ! "$(GOFMT_OUTPUT)" ]; then \
		echo "gofmt success"; \
	else \
		echo "gofmt failure. Please run:"; \
		echo "  gofmt -w -s $(GOFMT_OUTPUT)"; \
		exit 1; \
	fi

dist: cross-build
	cd dist && \
	$(DIST_DIRS) cp ../LICENSE {} \; && \
	$(DIST_DIRS) cp ../README.md {} \; && \
	$(DIST_DIRS) cp ../.config.json {} \; && \
	$(DIST_DIRS) tar -zcf $(NAME)-$(VERSION)-{}.tar.gz {} \; && \
	cd ..
