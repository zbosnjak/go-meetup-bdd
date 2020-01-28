GOARCH:=amd64
GOOS:=darwin
VERSION:=1.0.0
BUILD:=$(shell date +%s)
RELEASE:=1
PRODUCT0:=go-meetup-bdd
PRODUCT1:=bank-acc
REPO:=github.com/zbosnjak

################################################################################
# BUILD
################################################################################
.PHONY: build-deps
build-deps:
	go mod tidy

.PHONY: build
build: build-deps ## build the application
	GOOS=$(GOOS) GOARCH=$(GOARCH) $(GOENVS) go build -v -ldflags \
		"-X main.Version=$(VERSION) -X main.Build=$(BUILD) $(LDFLAGS)" \
		$(BUILDFLAGS) -o $(PRODUCT1) "$(REPO)/$(PRODUCT0)/cmd/$(PRODUCT1)"

.PHONY: clean
clean: ## clean up build PRODUCTs and rpms
	go clean ./...
	rm -f $(PRODUCT)
	rm -rf vendor
	rm -rf test/{results,reports}

.PHONY: vendor
vendor: ## create vendor directory of dependencies (used for docker images and for debugging)
	go mod vendor

################################################################################
# LINT
################################################################################

.PHONY: lint-deps
lint-deps: ## get linter for testing
	go get golang.org/x/lint/golint

.PHONY: lint
lint: lint-deps ## get linter for testing
	 golint -set_exit_status `(go list ./...)`

################################################################################
# TEST
################################################################################

.PHONY: tests-run
tests-run:  ## runs the automated acceptance test suite on the test environment
	source assets/envs/test && \
	bash scripts/test-suite.sh

.PHONY: tests-report
tests-report: ## generate html report from result.json files
	docker build -f build/docker/reporter.Dockerfile -t cucumber-html-reporter .
	docker run --rm \
		-v $(CURDIR)/tests:/test \
		cucumber-html-reporter

################################################################################
# HELP
################################################################################

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
