.DEFAULT_GOAL := all

.PHONY: all
all: test lint

# the TEST_FLAGS env var can be set to eg run only specific tests
# the coverage output can be open with
# `go tool cover -html=coverage.out`
TEST_COMMAND = go test -v -count=1 -race -cover -coverprofile=coverage.out $(TEST_FLAGS)

.PHONY: test
test:
	$(TEST_COMMAND)

# TODO wkpo!
.PHONY: lint
lint:
	golangci-lint run
