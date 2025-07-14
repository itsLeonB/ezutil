.PHONY: help lint test test-verbose test-coverage test-clean

help:
	@echo "Available commands:"
	@echo "  help           - Show this help message"
	@echo "  lint           - Run golangci-lint on the codebase"
	@echo "  test           - Run all tests"
	@echo "  test-verbose   - Run all tests with verbose output"
	@echo "  test-coverage  - Run all tests with coverage report"
	@echo "  test-clean     - Clean test cache and run tests"

lint:
	golangci-lint run ./...

test:
	@echo "Running all tests..."
	cd test && go test ./...

test-verbose:
	@echo "Running all tests with verbose output..."
	cd test && go test -v ./...

test-coverage:
	@echo "Running all tests with coverage report..."
	cd test && go test -v -cover ./...

test-clean:
	@echo "Cleaning test cache and running tests..."
	cd test && go clean -testcache && go test -v ./...
