default: build

build:
	cd cli/client; go build -v -o ../../bin/client-socks5

run: build
	./bin/socks5

clean:
	rm -rf bin

test:
	go test ./src/
