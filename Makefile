SHELL=/bin/bash
cgo = 1
version = 1.0.0
user= hrishikesh

.PHONY: build
build: build-arm64 build-amd64

.PHONY: build-arm64
build-arm64:
	@echo "👷 Building... build/hash-arm64-v${version}"
	@CGO_ENABLED=${cgo} GOOS=linux GOARCH=arm64 CC="zig cc -target aarch64-linux-musl" CXX="zig c++ -target aarch64-linux-musl" go build -o build/hash_arm64_v${version}
	@echo "✅ build/hash-arm64-v${version}"

.PHONY: build-amd64
build-amd64:
	@echo "👷 Building... build/hash-amd64-v${version}"
	@CGO_ENABLED=${cgo} GOOS=linux GOARCH=amd64 CC="zig cc -target x86_64-linux-musl" CXX="zig c++ -target x86_64-linux-musl" go build -o build/hash_amd64_v${version}
	@echo "✅ build/hash-amd64-v${version}"

.PHONY: clean
clean:
	rm -rf build/