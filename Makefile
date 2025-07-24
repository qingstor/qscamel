SHELL := /bin/bash

.PHONY: all check format　vet lint build install uninstall release clean test coverage

VERSION=$(shell cat ./constants/version.go | grep "Version\ =" | sed -e s/^.*\ //g | sed -e s/\"//g)

help:
	@echo "Please use \`make <target>\` where <target> is one of"
	@echo "  check      to format, vet and lint "
	@echo "  build      to create bin directory and build qscamel"
	@echo "  install    to install qscamel to /usr/local/bin/qscamel"
	@echo "  uninstall  to uninstall qscamel"
	@echo "  release    to release qscamel"
	@echo "  clean      to clean build and test files"
	@echo "  test       to run test"
	@echo "  coverage   to test with coverage"

check: format vet

format:
	@echo "go fmt"
	@go fmt ./...
	@echo "ok"

vet:
	@echo "go vet"
	@go vet ./...
	@echo "ok"

tidy:
	@echo "Tidy and check the go mod files"
	@go mod tidy && go mod verify
	@echo "Done"

build: tidy check
	@echo "build qscamel"
	@mkdir -p ./bin
	@CGO_ENABLED=0 go build -o ./bin/qscamel .
	@echo "ok"

install: build
	@echo "install qscamel to GOPATH"
	@cp ./bin/qscamel ${GOPATH}/bin/qscamel
	@echo "ok"

uninstall:
	@echo "delete /usr/local/bin/qscamel"
	@rm -f /usr/local/bin/qscamel
	@echo "ok"

release:
	@echo "release qscamel"
	@-rm ./release/*
	@mkdir -p ./release

	@echo "build for linux"
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags netgo -o ./bin/linux/qscamel_v${VERSION}_linux_amd64 .
	@tar -C ./bin/linux/ -czf ./release/qscamel_v${VERSION}_linux_amd64.tar.gz qscamel_v${VERSION}_linux_amd64

	@echo "build for linux arm64"
	@GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -tags netgo -o ./bin/linux/qscamel_v${VERSION}_linux_arm64 .
	@tar -C ./bin/linux/ -czf ./release/qscamel_v${VERSION}_linux_arm64.tar.gz qscamel_v${VERSION}_linux_arm64

	@echo "build for macOS"
	@GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -tags netgo -o ./bin/macos/qscamel_v${VERSION}_macos_amd64 .
	@tar -C ./bin/macos/ -czf ./release/qscamel_v${VERSION}_macos_amd64.tar.gz qscamel_v${VERSION}_macos_amd64

	@echo "build for macOS arm64"
	@GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -tags netgo -o ./bin/macos/qscamel_v${VERSION}_macos_arm64 .
	@tar -C ./bin/macos/ -czf ./release/qscamel_v${VERSION}_macos_arm64.tar.gz qscamel_v${VERSION}_macos_arm64

	@echo "build for windows"
	@GOOS=windows GOARCH=386 CGO_ENABLED=0 go build -tags netgo -o ./bin/windows/qscamel_v${VERSION}_windows_i386.exe .
	@tar -C ./bin/windows/ -czf ./release/qscamel_v${VERSION}_windows_i386.tar.gz qscamel_v${VERSION}_windows_i386.exe
	@GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -tags netgo -o ./bin/windows/qscamel_v${VERSION}_windows_amd64.exe .
	@tar -C ./bin/windows/ -czf ./release/qscamel_v${VERSION}_windows_amd64.tar.gz qscamel_v${VERSION}_windows_amd64.exe

	@echo "ok"

clean:
	@rm -rf ./bin
	@rm -rf ./release
	@rm -rf ./coverage

test:
	@echo "run test"
	@go test -v ./...
	@echo "ok"

coverage:
	@echo "run test with coverage"
	@go test -race -coverprofile=coverage.txt -covermode=atomic -v ./...
	@go tool cover -html="coverage.txt" -o "coverage.html"
	@echo "ok"
