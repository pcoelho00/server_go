BINARY_NAME=out

build:
	go build -o out

run: build
	./${BINARY_NAME}

clean:
	go clean
	rm ${BINARY_NAME}

test:
	go test -v ./...
