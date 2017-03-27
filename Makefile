default: build

build:
	go build -v -o ./bin/socks5

run: build
	./bin/socks5

clean:
	rm -rf bin

test:
	go test ./src/
