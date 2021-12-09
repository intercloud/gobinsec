BUILD_DIR = build

.DEFAULT_GOAL :=
default: clean fmt lint test integ

clean: # Clean generated files
	@rm -rf $(BUILD_DIR)

fmt: # Format Go code
	@go fmt ./...

lint: # Check Go code
	@golangci-lint run ./...

.PHONY: test
test: # Run unit tests
	@go test -cover ./...

.PHONY: build
build: # Build binary
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/ ./...

integ: build # Run integration test
	-@$(BUILD_DIR)/gobinsec test/binary > $(BUILD_DIR)/report.yml
	@test $? || (echo "ERROR should have exited with code 1" && exit 1)
	@cmp test/report.yml $(BUILD_DIR)/report.yml
	@$(BUILD_DIR)/gobinsec -verbose -config test/config.yml test/binary > $(BUILD_DIR)/report-config.yml
	@cmp test/report-config.yml $(BUILD_DIR)/report-config.yml

binaries: # Generate binaries
	@GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/gobinsec-linux-amd64 .
	@GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/gobinsec-darwin-amd64 .
	@GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/gobinsec-linux-arm64 .
