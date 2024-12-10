TEST_FLAGS ?= -p 8 -failfast -race -shuffle on

all:
	@echo "make <cmd>:"
	@echo ""
	@echo "commands:"
	@awk -F'[ :]' '/^#+/ {comment=$$0; gsub(/^#+[ ]*/, "", comment)} !/^(_|all:)/ && /^([A-Za-z_-]+):/ && !seen[$$1]++ {printf "  %-24s %s\n", $$1, (comment ? "- " comment : ""); comment=""} !/^#+/ {comment=""}' Makefile

test-clean:
	go clean -testcache

test: test-clean
	go test -run=$(TEST) $(TEST_FLAGS) -json ./... | tparse --all --follow

test-rerun: test-clean
	go run github.com/goware/rerun/cmd/rerun -watch ./ -run 'make test'

test-coverage:
	go test -run=$(TEST) $(TEST_FLAGS) -cover -coverprofile=coverage.out -json ./... | tparse --all --follow

test-coverage-inspect: test-coverage
	go tool cover -html=coverage.out

lint:
	golangci-lint run ./... --fix -c .golangci.yml

