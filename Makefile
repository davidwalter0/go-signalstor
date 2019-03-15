SHELL=/bin/bash
MAKEFILE_DIR := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
CURRENT_DIR := $(notdir $(patsubst %/,%,$(dir $(MAKEFILE_DIR))))
DIR=$(MAKEFILE_DIR)
PWD:=$(shell pwd)
PKG:=main
package_files:=$(filter-out %_test.go,$(wildcard *.go))
files:=$(filter-out examples/version.go %_test.go,$(wildcard examples/*.go))
VERSION_FLAG_INFO:=-X $(PKG).Version=$(shell cat .version) -X $(PKG).Build=$$(date -u +%Y.%m.%d.%H.%M.%S.%:::z) -X $(PKG).Commit=$$(git log --format=%h-%aI -n1)
targets:=$(patsubst %.go,bin/%,$(wildcard *.go))

all: $(patsubst examples/%.go,bin/%,$(files)) examples/version.go

bin/%: examples/%.go $(package_files)
	mkdir -p examples
	echo $(VERSION_FLAG_INFO)
# go build -tags netgo -ldflags "$(VERSION_FLAG_INFO)" -o $@ $< examples/version.go
	vgo build -tags netgo -ldflags "$(VERSION_FLAG_INFO)" -o $@ $< examples/version.go

.PHONY: test

test: 
# make -C test -f ../Makefile go-test
	cd test; go test -v 

go-test:
	cd test; go test -v 

install:
	go install

clean:
	rm -f $(MAKEFILE_DIR)/bin/*
