
server: deps
	go build -o bin/rapid

deps:
	go get -v ./...