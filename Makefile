#!/usr/bin/env make 
GO_BUILD     := go build -ldflags "-s -w" -a -installsuffix cgo
GO_TEST      := go test ./... -v -cover

build: ## builds the binary
	@$(GO_BUILD) -o api .

.PHONY: test
test: ## runs all tests
	@$(GO_TEST)

deploy: ## deploy kube-prometheus-stack and the API 
	./scripts/deploy-kps.sh
	./scripts/deploy-api.sh