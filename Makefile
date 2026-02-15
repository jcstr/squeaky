APP      := squeaky
VERSION  := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS  := -s -w -X github.com/jcstr/squeaky/cmd.version=$(VERSION)

.PHONY: build install test lint clean

build:
	go build -ldflags "$(LDFLAGS)" -o $(APP) .

install: build
	install -Dm755 $(APP) /usr/local/bin/$(APP)

test:
	go test ./... -v -race

lint:
	golangci-lint run ./...

clean:
	rm -f $(APP)
