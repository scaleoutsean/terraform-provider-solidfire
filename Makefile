-include $(HOME)/.tf-elementsw-devrc.mk

TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)

default: build

build: fmtcheck
	go install

test: fmtcheck
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

testacc-qos-policy: fmtcheck
	TF_ACC=1 go test ./elementsw -v -run TestAccElementswQoSPolicy

testacc-volume: fmtcheck
	TF_ACC=1 go test ./elementsw -v -run TestAccElementswVolume

testacc-account: fmtcheck
	TF_ACC=1 go test ./elementsw -v -run TestAccount_

testacc-initiator: fmtcheck
	TF_ACC=1 go test ./elementsw -v -run TestAccElementswInitiator

testacc-pairing: fmtcheck
	TF_ACC=1 go test ./elementsw -v -run TestAccElementsw.*Pairing

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

vendor-status:
	@govendor status

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./aws"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

.PHONY: build test testacc testacc-qos-policy testacc-volume testacc-account testacc-initiator testacc-pairing vet fmt fmtcheck errcheck vendor-status test-compile
