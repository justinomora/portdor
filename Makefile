BINARY=portdor
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

.PHONY: build test clean install

build:
	go build -ldflags "-X main.Version=$(VERSION)" -o $(BINARY) ./cmd/portdor

test:
	go test ./...

install: build
	cp $(BINARY) /usr/local/bin/

clean:
	rm -f $(BINARY)
