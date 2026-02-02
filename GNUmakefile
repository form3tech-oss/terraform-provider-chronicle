SHELL := /bin/bash

PROJECT_NAME := $(shell basename $(CURDIR))
empty:=
prefix:=terraform-provider-
PROVIDER_NAME := $(subst $(prefix),$(empty),$(PROJECT_NAME))

OS := $(shell go env GOHOSTOS)
ARCH := $(shell go env GOHOSTARCH)
LOCAL_PLUGIN_DIR := ~/.terraform.d/plugins/github.com/form3tech-oss/$(PROVIDER_NAME)/0.0.1/$(OS)_$(ARCH)


GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)

default: lint build test

test: fmtcheck
	go test -v . ./chronicle

testacc: fmtcheck
	@echo "WARNING: Running both V1 and V2 feed tests together. This may fail if your Chronicle SIEM"
	@echo "         has both V1 and V2 feeds active. Consider using 'make testacc-v1' or 'make testacc-v2' instead."
	@echo ""
	TF_ACC=1  go test -v ./chronicle -timeout 120m  -parallel 1

testacc-v1: fmtcheck
	TF_ACC=1  go test -v ./chronicle -run='TestAcc' -skip='V2|EventDriven' -timeout 120m  -parallel 1

testacc-v2: fmtcheck
	TF_ACC=1  go test -v ./chronicle -run='TestAcc.*(V2|EventDriven)' -timeout 120m  -parallel 1	

build:
	@go build -mod=vendor -o $(PROJECT_NAME)
	@echo "Build succeeded"

install: lint build
	@mkdir -p $(LOCAL_PLUGIN_DIR)
	@cp $(PROJECT_NAME) $(LOCAL_PLUGIN_DIR)/
	@echo "Install succeeded"

clean:
	go clean -testcache

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v /vendor | grep -v /tools) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

lint: fmtcheck vet

docs:
	tfplugindocs generate

vendor:
	@go mod tidy && go mod vendor && go mod verify

.PHONY: build install lint test clean testacc testacc-v1 testacc-v2 vet fmt fmtcheck docs vendor
