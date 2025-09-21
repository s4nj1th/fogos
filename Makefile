BINARY_NAME = fogos
MAIN_PATH   = ./cmd/fogos
INSTALL_PATH= /usr/local/bin
BUILD_DIR   = ./bin

.PHONY: build clean install uninstall test run-example fmt vet

build:
	@mkdir -p $(BUILD_DIR)
	go build -buildvcs=false -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

build-all:
	@mkdir -p $(BUILD_DIR)
	GOOS=linux  GOARCH=amd64 go build -buildvcs=false -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64  $(MAIN_PATH)
	GOOS=darwin GOARCH=amd64 go build -buildvcs=false -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=arm64 go build -buildvcs=false -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	GOOS=windows GOARCH=amd64 go build -buildvcs=false -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)

clean:
	go clean
	rm -rf $(BUILD_DIR)

install: build
	cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_PATH)/
	chmod +x $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "fogos installed to $(INSTALL_PATH)/$(BINARY_NAME)"

uninstall:
	rm -f $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "fogos uninstalled"

test:
	go test -v ./...

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

fmt:
	go fmt ./...

vet:
	go vet ./...

check: fmt vet test

dev-run: build
	@echo "Building fogos..."
	@echo "Example usage (requires root):"
	@echo "  $(BUILD_DIR)/fogos block youtube.com"
	@echo "  $(BUILD_DIR)/fogos status youtube.com"
	@echo "  $(BUILD_DIR)/fogos list"
	@echo "  $(BUILD_DIR)/fogos unblock youtube.com"

help:
	@echo "Available targets:"
	@echo "  build         - Build the binary"
	@echo "  build-all     - Build for multiple platforms"
	@echo "  clean         - Clean build artifacts"
	@echo "  install       - Install to system (run as root)"
	@echo "  uninstall     - Uninstall from system (run as root)"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  fmt           - Format code"
	@echo "  vet           - Run go vet"
	@echo "  check         - Run fmt, vet, and test"
	@echo "  dev-run       - Build and show usage examples"
	@echo "  help          - Show this help"
