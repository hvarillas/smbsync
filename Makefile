.PHONY: build build-windows build-linux test clean run

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
	go test ./...

clean:
	rm -f smbsync smbsync-* *.exe *.log

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
