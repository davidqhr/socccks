default: build

build:
	cd cli/client; go build -v -o ../../bin/socccks-client
	cd cli/server; go build -v -o ../../bin/socccks-server

run: build
	./bin/socccks

clean:
	rm -rf bin

test:
	go test -test.bench=".*" ./client
	go test -test.bench=".*" ./server
	go test -test.bench=".*" ./utils
