BUILD_DIR = build

.DEFAULT_GOAL :=
default: fmt lint test run

clean: # Clean generated files
	@rm -rf $(BUILD_DIR)

fmt: # Format Go code
	@go fmt *.go

lint: # Check Go code
	@golangci-lint run ./...

test: # Run unit tests
	@go test -cover ./...

.PHONY: build
build: # Build binary
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/ ./...

run: build # Run tool
	@$(BUILD_DIR)/gobinsec ../nancy-test/build/nancy-test
	@#$(BUILD_DIR)/gobinsec -verbose ../user-orga/build/user-orga
