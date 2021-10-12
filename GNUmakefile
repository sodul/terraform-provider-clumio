#
# Copyright 2021. Clumio, Inc.
#

VERSION=0.1.1
OS_ARCH=darwin_amd64

CLUMIO_PROVIDER_DIR=~/.terraform.d/plugins/clumio.com/providers/clumio/${VERSION}/${OS_ARCH}

default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

install:
	go mod vendor
	mkdir -p ${CLUMIO_PROVIDER_DIR}
	go build -o ${CLUMIO_PROVIDER_DIR}/terraform-provider-clumio_v${VERSION}
