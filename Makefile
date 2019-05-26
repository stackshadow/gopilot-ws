
deps:
	GOPATH=${HOME}/go:${PWD}/../gopilot-lib:${PWD} \
	go get -d -v ./src

build:
	GOPATH=${HOME}/go:${PWD}/../gopilot-lib:${PWD} \
	go build -o gopilot-ws ./src

clean:
	rm gopilot-ws