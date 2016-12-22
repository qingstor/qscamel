SHELL := /bin/bash

.PHONY: all check format　vet lint build install uninstall release clean test coverage

VERSION=$(shell cat ./metadata/version.go | grep "Version\ =" | sed -e s/^.*\ //g | sed -e s/\"//g)
DIRS_TO_CHECK=$(shell ls -d */ | grep -vE "vendor|test")
PKGS_TO_CHECK=$(shell go list ./... | grep -v "/vendor/")
SUPPORTED_OS=linux darwin windows

help:
	@echo "Please use \`make <target>\` where <target> is one of"
	@echo "  all        to check, build and test qscamel"
	@echo "  check      to format, vet and lint "
	@echo "  build      to create bin directory and build qscamel"
	@echo "  install    to install qscamel to /usr/local/bin/qscamel"
	@echo "  uninstall  to uninstall qscamel"
	@echo "  release    to release qscamel"
	@echo "  clean      to clean build and test files"
	@echo "  test       to run test"
	@echo "  coverage   to test with coverage"

all: check build test

check: format vet lint

format:
	@echo "go fmt, skipping vendor packages"
	@for pkg in ${PKGS_TO_CHECK}; do go fmt $${pkg}; done;
	@echo "ok"

vet:
	@echo "go tool vet, skipping vendor packages"
	@go tool vet -all ${DIRS_TO_CHECK}
	@echo "ok"

lint:
	@echo "golint, skipping vendor packages"
	@lint=$$(for pkg in ${PKGS_TO_CHECK}; do golint $${pkg}; done); \
	 lint=$$(echo "$${lint}"); \
	 if [[ -n $${lint} ]]; then echo "$${lint}"; exit 1; fi
	@echo "ok"

build:
	@echo "build qscamel"
	@mkdir -p ./bin
	@go build -o ./bin/qscamel .
	@echo "ok"

install: build
	@echo "install qscamel to /usr/local/bin/qscamel"
	@cp ./bin/qscamel /usr/local/bin/qscamel
	@echo "ok"

uninstall:
	@echo "delete /usr/local/bin/qscamel"
	@rm -f /usr/local/bin/qscamel
	@echo "ok"

release:
	@echo "release qscamel"
	@mkdir -p ./release
	@for os in ${SUPPORTED_OS}; do \
		echo "for $${os}"; \
		mkdir -p ./bin/$${os}; \
		SUPPORTED_OS=$${os} GOARCH=386 go build -o ./bin/$${os}/qscamel_v${VERSION}_$${os}_32 .; \
	    tar -C ./bin/$${os}/ -czf ./release/qscamel_v${VERSION}_$${os}_32.tar.gz qscamel_v${VERSION}_$${os}_32; \
	    SUPPORTED_OS=$${os} GOARCH=amd64 go build -o ./bin/$${os}/qscamel_v${VERSION}_$${os}_64 .; \
	    tar -C ./bin/$${os}/ -czf ./release/qscamel_v${VERSION}_$${os}_64.tar.gz qscamel_v${VERSION}_$${os}_64; \
	done
	@echo "ok"

clean:
	@rm -rf ./bin
	@rm -rf ./release
	@rm -rf ./coverage

test:
	@echo "run test"
	@go test -v ${PKGS_TO_CHECK}
	@echo "ok"

coverage:
	@echo "run test with coverage"
	@for pkg in ${PKGS_TO_CHECK}; do \
		output="coverage$${pkg#github.com/yunify/qscamel}"; \
		mkdir -p $${output}; \
		go test -v -cover -coverprofile="$${output}/profile.out" $${pkg}; \
		if [[ -e "$${output}/profile.out" ]]; then \
			go tool cover -html="$${output}/profile.out" -o "$${output}/profile.html"; \
		fi; \
	done
	@echo "ok"
