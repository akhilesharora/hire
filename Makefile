PHONY: help

MIN_COVERAGE=70

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Compile project
	go build -o build/hire

rebuild: ## Clean & Reinstall dependencies
	rm -r vendor
	rm go.sum
	go build -o build/hire

gofmt: ## Format
	./scripts/gofmt.sh .

gometalinter: ## Linter
	./scripts/gometalinter.sh .

test: ## Run all tests
	go test github.com/akhilesharora/hire/internal/config
	go test github.com/akhilesharora/hire/internal

coverage: ## Run tests and generate coverage files per package
	mkdir .coverage 2> /dev/null || true
	rm -rf .coverage/*.out || true
	go test github.com/akhilesharora/hire/internal/config -coverprofile=.coverage/config.out
	go test github.com/akhilesharora/hire/internal -coverprofile=.coverage/internal.out

coverage-total: coverage-concat ## Total coverage of all packages
	./scripts/coverage.sh $(MIN_COVERAGE)

coverage-html: coverage-concat ## Open html results
	go tool cover -html=.coverage/all.out

coverage-concat: ## Concat cover files into one
	rm -rf .coverage/all.out 2> /dev/null || true
	echo "mode: set" > .coverage/all.out
	cat .coverage/*.out | grep -v "mode: set" >> .coverage/all.out

docker-build:
	go mod vendor
	docker build -t hire --no-cache -f Dockerfile .
