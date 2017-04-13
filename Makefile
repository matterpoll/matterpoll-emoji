NAME	 := matterpoll-emoji
VERSION  := 0.0.1
REVISION := $(shell git rev-parse --short HEAD)

LDFLAGS := -ldflags="-s -w -X \"main.Version=$(VERSION)\" -X \"main.Revision=$(REVISION)\" -extldflags \"-static\""
DIST_DIRS := find * -type d -exec

.PHONY: glide deps clean test cross-build dist

all: dist test

glide:
ifeq ($(shell command -v glide 2> /dev/null),)
	curl https://glide.sh/get | sh
endif

deps: glide
	glide install

cross-build: deps
	for os in darwin linux windows; do \
		for arch in amd64 386; do \
			GOOS=$$os GOARCH=$$arch CGO_ENABLED=0 go build -a -tags netgo -installsuffix netgo $(LDFLAGS) -o dist/$$os-$$arch/$(NAME); \
		done; \
	done 

clean:
	rm -rf bin/*
	rm -rf vendor/*

test:
	go test -cover -v `glide novendor`

dist: cross-build
	cd dist && \
	$(DIST_DIRS) cp ../LICENSE {} \; && \
	$(DIST_DIRS) cp ../README.md {} \; && \
	$(DIST_DIRS) tar -zcf $(NAME)-$(VERSION)-{}.tar.gz {} \; && \
	$(DIST_DIRS) zip -r $(NAME)-$(VERSION)-{}.zip {} \; && \
	cd ..
