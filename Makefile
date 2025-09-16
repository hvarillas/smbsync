.PHONY: build build-windows build-linux test test-verbose test-coverage test-race test-bench clean run

build:
	@echo "Building smbsync..."
	go build -ldflags "-s -w" -trimpath -o smbsync ./cmd/smbsync

build-windows:
	@echo "Building for Windows..."
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-s -w" -trimpath -o smbsync-windows-amd64.exe ./cmd/smbsync

build-linux:
	@echo "Building for Linux..."
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-s -w" -trimpath -o smbsync-linux-amd64 ./cmd/smbsync

build-windows-upx: build-windows
	@echo "Compressing with UPX..."
	upx --best --lzma smbsync-windows-amd64.exe

test:
	@echo "Running unit tests..."
	go test ./...

test-verbose:
	@echo "Running unit tests with verbose output..."
	go test -v ./...

test-coverage:
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-race:
	@echo "Running tests with race detection..."
	go test -race ./...

test-bench:
	@echo "Running benchmarks..."
	go test -bench=. ./...

clean:
	rm -f smbsync smbsync-* *.exe *.log coverage.out coverage.html

run:
	@echo "Make sure to copy .env.example to .env and configure it first"
	go run ./cmd/smbsync --help

encrypt-pass:
	@echo "Usage: make encrypt-pass PASS=your_password"
	go run ./cmd/smbsync --generate-encrypted --pass "$(PASS)"

fmt:
	go fmt ./...

vet:
	go vet ./...

mod-tidy:
	go mod tidy
