SHELL=/bin/bash
cgo = 1
version = 1.0.0

.PHONY: build
build: build-linux-arm64 build-linux-amd64 build-windows-amd64 build-windows-arm64

.PHONY: build-linux-arm64
build-linux-arm64:
	@echo "👷 Building... build/hash-linux-arm64-v${version}"
	@CGO_ENABLED=${cgo} GOOS=linux GOARCH=arm64 CC="zig cc -target aarch64-linux-musl" CXX="zig c++ -target aarch64-linux-musl" go build -o build/hash_linux_arm64_v${version}
	@echo "✅ build/hash-linux-arm64-v${version}"

.PHONY: build-linux-amd64
build-linux-amd64:
	@echo "👷 Building... build/hash-linux-amd64-v${version}"
	@CGO_ENABLED=${cgo} GOOS=linux GOARCH=amd64 CC="zig cc -target x86_64-linux-musl" CXX="zig c++ -target x86_64-linux-musl" go build -o build/hash_linux_amd64_v${version}
	@echo "✅ build/hash-linux-amd64-v${version}"

.PHONY: build-windows-amd64
build-windows-amd64:
	@echo "👷 Building... build/hash_windows_amd64_v${version}.exe"
	@CGO_ENABLED=${cgo} GOOS=windows GOARCH=amd64 CC="zig cc -target x86_64-windows-gnu" CXX="zig c++ -target x86_64-windows-gnu"  go build -o build/hash_windows_amd64_v${version}.exe
	@echo "✅ build/hash_windows_amd64_v${version}.exe"

.PHONY: build-windows-arm64
build-windows-arm64:
	@echo "👷 Building... build/hash_windows_arm64_v${version}.exe"
	@CGO_ENABLED=${cgo} GOOS=windows GOARCH=arm64 CC="zig cc -target aarch64-windows-gnu" CXX="zig c++ -target aarch64-windows-gnu"  go build -o build/hash_windows_arm64_v${version}.exe
	@echo "✅ build/hash_windows_arm64_v${version}.exe"

.PHONY: clean
clean:
	rm -rf build/