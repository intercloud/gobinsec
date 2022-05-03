BUILD_DIR   = build
VERSION     = "UNKNOWN"
GOOSARCH    = $(shell go tool dist list | grep -v android)
MAIN_BRANCH = publish-release

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

.PHONY: install
install: # Build and install tool
	@go install .

integ: build # Run integration test
	-@$(BUILD_DIR)/gobinsec test/binary > $(BUILD_DIR)/report.yml
	@test $? || (echo "ERROR should have exited with code 1" && exit 1)
	@cmp test/report.yml $(BUILD_DIR)/report.yml
	@cat test/config.yml | envsubst > $(BUILD_DIR)/config.yml
	@$(BUILD_DIR)/gobinsec -verbose -config $(BUILD_DIR)/config.yml test/binary > $(BUILD_DIR)/report-config.yml
	@cmp test/report-config.yml $(BUILD_DIR)/report-config.yml

binaries: # Generate binaries
	@mkdir -p $(BUILD_DIR)/bin
	@gox -ldflags "-X main.Version=$(VERSION) -s -f" -osarch '$(GOOSARCH)' -output=$(BUILD_DIR)/bin/{{.Dir}}-{{.OS}}-{{.Arch}} $(GOPACKAGE)
	@cp install $(BUILD_DIR)/bin/

check: # Check release prerequisites
	@if [ "$(VERSION)" = "UNKNOWN" ]; then \
		echo "ERROR you must pass VERSION=X.Y.Z on command line"; \
		exit 1; \
	fi
	@if [ "$(TITLE)" = "EMPTY" ]; then \
		echo 'ERROR you must pass TITLE="..." on command line'; \
		exit 1; \
	fi
	@if [ "$$GITHUB_USER" = "" ]; then \
		echo "GITHUB_USER must be defined in your environment"; \
		exit 1; \
	fi
	@if [ "$$GITHUB_TOKEN" = "" ]; then \
		echo "GITHUB_TOKEN must be defined in your environment"; \
		exit 1; \
	fi
	@git diff-index --quiet HEAD -- || (echo "ERROR There are uncommitted changes" && exit 1)
	@test `git rev-parse --abbrev-ref HEAD` = "$(MAIN_BRANCH)" || (echo "ERROR You are not on branch $(MAIN_BRANCH)" && exit 1)

tag: # Create release tag
	@git tag -a $(VERSION) -m "Release $(VERSION)"
	@git push origin --tags

upload: # Publish release on github
	@echo "Creating release $(VERSION)"
	@read -p "Title: " title; \
	description=`git log --pretty=format:"%h %s"`; \
	github-release release \
		--user intercloud \
		--repo gobinsec \
		--tag "$(VERSION)" \
		--name "$$title" \
		--description "$$description"
	@sleep 5
	@for file in $(BUILD_DIR)/bin/*; do \
		echo "Uploading $$file..."; \
		github-release upload \
			--user intercloud \
			--repo gobinsec \
			--tag "$(VERSION)" \
			--name `basename $$file` \
			--file $$file; \
	done

release: check clean lint test integ binaries tag upload # Perform release (must pass VERSION=X.Y.Z on command line)
