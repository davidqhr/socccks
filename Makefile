default: build

build:
	go get -d ./...
	cd cli/socccks-client; go build -v -o ../../bin/socccks-client
	cd cli/socccks-server; go build -v -o ../../bin/socccks-server

clean:
	rm -rf bin

test:
	go test -test.bench=".*" ./client
	go test -test.bench=".*" ./server
	go test -test.bench=".*" ./utils
