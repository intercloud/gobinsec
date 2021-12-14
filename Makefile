BUILD_DIR = build
VERSION   = "UNKNOWN"

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
	@go build -ldflags "-X main.Version=$(VERSION) -s -f" -o $(BUILD_DIR)/ ./...

install: # Build and install tool
	@go install .

integ: build # Run integration test
	-@$(BUILD_DIR)/gobinsec test/binary > $(BUILD_DIR)/report.yml
	@test $? || (echo "ERROR should have exited with code 1" && exit 1)
	@cmp test/report.yml $(BUILD_DIR)/report.yml
	@$(BUILD_DIR)/gobinsec -verbose -config test/config.yml test/binary > $(BUILD_DIR)/report-config.yml
	@cmp test/report-config.yml $(BUILD_DIR)/report-config.yml

binaries: # Generate binaries
	@GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Version=$(VERSION) -s -f" -o $(BUILD_DIR)/gobinsec-linux-amd64 .
	@GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.Version=$(VERSION) -s -f" -o $(BUILD_DIR)/gobinsec-darwin-amd64 .
	@GOOS=darwin GOARCH=arm64 go build -ldflags "-X main.Version=$(VERSION) -s -f" -o $(BUILD_DIR)/gobinsec-darwin-arm64 .

version: # Check that VERSION=X.Y.Z was passed on command line
	@if [ "$(VERSION)" = "UNKNOWN" ]; then \
		echo "ERROR you must pass VERSION=X.Y.Z on command line"; \
		exit 1; \
	fi

release: version clean lint test integ binaries # Perform release (must pass VERSION=X.Y.Z on command line)
	@git tag $(VERSION)
