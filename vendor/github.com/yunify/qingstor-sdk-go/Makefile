SHELL := /bin/bash

.PHONY: all check vet lint update generate test build unit release clean

PREFIX=qingstor-sdk-go
VERSION=$(shell cat version.go | grep "Version\ =" | sed -e s/^.*\ //g | sed -e s/\"//g)
DIRS_TO_CHECK=$(shell ls -d */ | grep -vE "vendor|test")
PKGS_TO_CHECK=$(shell go list ./... | grep -v "/vendor/")
PKGS_TO_RELEASE=$(shell go list ./... | grep -vE "/vendor/|/test")
FILES_TO_RELEASE=$(shell find . -name "*.go" | grep -vE "/vendor/|/test|.*_test.go")
FILES_TO_RELEASE_WITH_VENDOR=$(shell find . -name "*.go" | grep -vE "/test|.*_test.go")

help:
	@echo "Please use \`make <target>\` where <target> is one of"
	@echo "  all               to check, build, test and release this SDK"
	@echo "  check             to vet and lint the SDK"
	@echo "  update            to update git submodules"
	@echo "  generate          to generate service code"
	@echo "  test              to run service test"
	@echo "  build             to build the SDK"
	@echo "  unit              to run all sort of unit tests except runtime"
	@echo "  unit-test         to run unit test"
	@echo "  unit-benchmark    to run unit test with benchmark"
	@echo "  unit-coverage     to run unit test with coverage"
	@echo "  unit-race         to run unit test with race"
	@echo "  unit-runtime      to run test with go1.7, go1.6, go 1.5 in docker"
	@echo "  release           to build and release current version"
	@echo "  release-source    to pack the source code"
	@echo "  release-headers   to build and pack the headers source code for go 1.7"
	@echo "  release-binary    to build the static binary for go 1.7"
	@echo "  clean             to clean the coverage files"

all: check build unit release

check: vet lint

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

update:
	git submodule update --remote
	@echo "ok"

generate:
	@if [[ ! -f "$$(which snips)" ]]; then \
		echo "ERROR: Command \"snips\" not found."; \
	fi
	snips \
		--service=qingstor --service-api-version=latest \
		--spec="./specs" --template="./template" --output="./service"
	gofmt -w .
	@echo "ok"

test:
	pushd "./test"; go run *.go; popd
	@echo "ok"

build:
	@echo "build the SDK"
	GOOS=linux GOARCH=amd64 go build ${PKGS_TO_CHECK}
	GOOS=darwin GOARCH=amd64 go build ${PKGS_TO_CHECK}
	GOOS=windows GOARCH=amd64 go build ${PKGS_TO_CHECK}
	@echo "ok"

unit: unit-test unit-benchmark unit-coverage unit-race

unit-test:
	@echo "run unit test"
	go test -v ${PKGS_TO_CHECK}
	@echo "ok"

unit-benchmark:
	@echo "run unit test with benchmark"
	go test -v -bench=. ${PKGS_TO_CHECK}
	@echo "ok"

unit-coverage:
	@echo "run unit test with coverage"
	for pkg in ${PKGS_TO_CHECK}; do \
		output="coverage$${pkg#github.com/yunify/qingstor-sdk-go}"; \
		mkdir -p $${output}; \
		go test -v -cover -coverprofile="$${output}/profile.out" $${pkg}; \
		if [[ -e "$${output}/profile.out" ]]; then \
			go tool cover -html="$${output}/profile.out" -o "$${output}/profile.html"; \
		fi; \
	done
	@echo "ok"

unit-race:
	@echo "run unit test with race"
	go test -v -race -cpu=1,2,4 ${PKGS_TO_CHECK}
	@echo "ok"

unit-runtime: unit-runtime-go-1.7 unit-runtime-go-1.6 unit-runtime-go-1.5

export define DOCKERFILE_GO_1_7
FROM golang:1.7

ADD . /go/src/github.com/yunify/qingstor-sdk-go
WORKDIR /go/src/github.com/yunify/qingstor-sdk-go

CMD ["make", "build", "unit"]
endef

unit-runtime-go-1.7:
	@echo "run test in go 1.7"
	echo "$${DOCKERFILE_GO_1_7}" > "dockerfile_go_1.7"
	docker build -f "./dockerfile_go_1.7" -t "${PREFIX}:go-1.7" .
	rm -f "./dockerfile_go_1.7"
	docker run --name "${PREFIX}-go-1.7-unit" -t "${PREFIX}:go-1.7"
	docker rm "${PREFIX}-go-1.7-unit"
	docker rmi "${PREFIX}:go-1.7"
	@echo "ok"

export define DOCKERFILE_GO_1_6
FROM golang:1.6

ADD . /go/src/github.com/yunify/qingstor-sdk-go
WORKDIR /go/src/github.com/yunify/qingstor-sdk-go

CMD ["make", "build", "unit"]
endef

unit-runtime-go-1.6:
	@echo "run test in go 1.6"
	echo "$${DOCKERFILE_GO_1_6}" > "dockerfile_go_1.6"
	docker build -f "./dockerfile_go_1.6" -t "${PREFIX}:go-1.6" .
	rm -f "./dockerfile_go_1.6"
	docker run --name "${PREFIX}-go-1.6-unit" -t "${PREFIX}:go-1.6"
	docker rm "${PREFIX}-go-1.6-unit"
	docker rmi "${PREFIX}:go-1.6"
	@echo "ok"

export define DOCKERFILE_GO_1_5
FROM golang:1.5
ENV GO15VENDOREXPERIMENT="1"

ADD . /go/src/github.com/yunify/qingstor-sdk-go
WORKDIR /go/src/github.com/yunify/qingstor-sdk-go

CMD ["make", "build", "unit"]
endef

unit-runtime-go-1.5:
	@echo "run test in go 1.5"
	echo "$${DOCKERFILE_GO_1_5}" > "dockerfile_go_1.5"
	docker build -f "dockerfile_go_1.5" -t "${PREFIX}:go-1.5" .
	rm -f "dockerfile_go_1.5"
	docker run --name "${PREFIX}-go-1.5-unit" -t "${PREFIX}:go-1.5"
	docker rm "${PREFIX}-go-1.5-unit"
	docker rmi "${PREFIX}:go-1.5"
	@echo "ok"

release: release-source release-source-with-vendor release-headers release-binary

release-source:
	@echo "pack the source code"
	mkdir -p "release"
	zip -FS "release/${PREFIX}-source-v${VERSION}.zip" ${FILES_TO_RELEASE}
	@echo "ok"

release-source-with-vendor:
	@echo "pack the source code"
	mkdir -p "release"
	zip -FS "release/${PREFIX}-source-with-vendor-v${VERSION}.zip" ${FILES_TO_RELEASE_WITH_VENDOR}
	@echo "ok"

release-headers: release-headers-go-1.7

release-headers-go-1.7:
	@echo "build and pack the headers source code for go 1.7"
	mkdir -p "release"
	mkdir -p "/tmp/${PREFIX}-headers/"
	for file in ${FILES_TO_RELEASE}; do \
		filepath="/tmp/${PREFIX}-headers/$$(dirname $${file})/binary.go"; \
		mkdir -p "$$(dirname $${filepath})"; \
		package_line=$$(cat "$${file}" | grep -E "^package"); \
		echo -ne "//go:binary-only-package\n\n" > "$${filepath}"; \
		echo -ne "$${package_line}" >> "$${filepath}"; \
	done
	pushd "/tmp/${PREFIX}-headers/"; \
	zip -r "/tmp/${PREFIX}-headers-v${VERSION}-go-1.7.zip" .; \
	popd
	cp "/tmp/${PREFIX}-headers-v${VERSION}-go-1.7.zip" "release/"
	rm -f "/tmp/${PREFIX}-headers-v${VERSION}-go-1.7.zip"
	rm -rf "/tmp/${PREFIX}-headers"
	@echo "ok"

release-binary: release-binary-go-1.7

release-binary-go-1.7:
	@echo "build the static binary for go 1.7"
	mkdir -p "release"
	for pkg in ${PKGS_TO_RELEASE}; do \
		GOOS=linux GOARCH=amd64 go install $${pkg}; \
		GOOS=darwin GOARCH=amd64 go install $${pkg}; \
		GOOS=windows GOARCH=amd64 go install $${pkg}; \
	done
	cross=(linux_amd64 darwin_amd64 windows_amd64); \
	for os_arch in $${cross[@]}; do \
		MAIN_GOPATH=$$(echo "${GOPATH}" | awk '{split($$1,p,":"); print(p[1])}'); \
		pushd "$${MAIN_GOPATH}/pkg/$${os_arch}/github.com/yunify/qingstor-sdk-go"; \
		zip -r "/tmp/${PREFIX}-binary-v${VERSION}-$${os_arch}-go-1.7.zip" .; \
		popd; \
		cp "/tmp/${PREFIX}-binary-v${VERSION}-$${os_arch}-go-1.7.zip" "release/"; \
		rm -f "/tmp/${PREFIX}-binary-v${VERSION}-$${os_arch}-go-1.7.zip"; \
	done
	@echo "ok"

clean:
	rm -rf $${PWD}/coverage
	@echo "ok"
