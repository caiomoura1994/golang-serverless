.PHONY: build clean deploy

build:
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/hello hello/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/world world/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/listClients listClients/main.go

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose

start: build
	sls offline --useDocker start --host 0.0.0.0  --verbose