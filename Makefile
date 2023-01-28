.PHONY: build clean deploy

build:
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/internal/listClients internal/listClients/main.go

clean:
	rm -rf ./bin

deploy:
	clean 
	build
	sls deploy --verbose --stage $(stage)

start: build
	sls offline --useDocker start --host 0.0.0.0  --verbose