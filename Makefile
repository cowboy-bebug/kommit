OUTPUT := kommit
VERSION := $(shell git describe --tags --always | sed 's/^v//')
COMMIT := $(shell git rev-parse --short HEAD)
DATE := $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
LDFLAGS := "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

build:
	go build -ldflags $(LDFLAGS) -o $(OUTPUT)

build-platform:
	mkdir -p dist
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags $(LDFLAGS) -o dist/$(OUTPUT)

release: release-darwin-arm64 release-darwin-amd64 release-linux-amd64 release-linux-arm64

release-darwin-arm64:
	$(MAKE) build-platform GOOS=darwin GOARCH=arm64
	tar -czf dist/$(OUTPUT)_v$(VERSION)_darwin_arm64.tar.gz -C dist $(OUTPUT)
	shasum -a 256 dist/$(OUTPUT)_v$(VERSION)_darwin_arm64.tar.gz > dist/$(OUTPUT)_v$(VERSION)_darwin_arm64.tar.gz.sha256

release-darwin-amd64:
	$(MAKE) build-platform GOOS=darwin GOARCH=amd64
	tar -czf dist/$(OUTPUT)_v$(VERSION)_darwin_amd64.tar.gz -C dist $(OUTPUT)
	shasum -a 256 dist/$(OUTPUT)_v$(VERSION)_darwin_amd64.tar.gz > dist/$(OUTPUT)_v$(VERSION)_darwin_amd64.tar.gz.sha256

release-linux-amd64:
	$(MAKE) build-platform GOOS=linux GOARCH=amd64
	tar -czf dist/$(OUTPUT)_v$(VERSION)_linux_amd64.tar.gz -C dist $(OUTPUT)
	shasum -a 256 dist/$(OUTPUT)_v$(VERSION)_linux_amd64.tar.gz > dist/$(OUTPUT)_v$(VERSION)_linux_amd64.tar.gz.sha256

release-linux-arm64:
	$(MAKE) build-platform GOOS=linux GOARCH=arm64
	tar -czf dist/$(OUTPUT)_v$(VERSION)_linux_arm64.tar.gz -C dist $(OUTPUT)
	shasum -a 256 dist/$(OUTPUT)_v$(VERSION)_linux_arm64.tar.gz > dist/$(OUTPUT)_v$(VERSION)_linux_arm64.tar.gz.sha256

clean:
	rm -rf dist

.PHONY: build build-platform release release-darwin-arm64 release-darwin-amd64 release-linux-amd64 release-linux-arm64 clean
