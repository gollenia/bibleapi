BINARY_NAME=bibleapi

.PHONY: setup run build build-linux fmt tidy clean

setup:
	go mod tidy
	go get github.com/gofiber/fiber/v2
	go get github.com/ncruces/go-sqlite3

run:
	go run .

build:
	go build -o $(BINARY_NAME) .

build-linux:
	GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o $(BINARY_NAME)_linux .

fmt:
	go fmt ./...

clean:
	rm -f $(BINARY_NAME) $(BINARY_NAME)_linux