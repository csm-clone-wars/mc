APP_NAME := mc

.PHONY: all build test clean run

all: build

build:
	go build -o $(APP_NAME) .

test:
	go test -v ./...

clean:
	rm -f $(APP_NAME)
	rm -rf include

run: build
	./$(APP_NAME)
