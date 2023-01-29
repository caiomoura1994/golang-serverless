.PHONY: build clean deploy

build:
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/internal/listClients internal/listClients/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/internal/upSertClient internal/upSertClient/main.go

clean:
	rm -rf ./bin

pre-build: clean build

deploy-qa: pre-build
	sls deploy --verbose --stage qa

deploy-prod: pre-build
	sls deploy --verbose --stage prod

deploy-dev: pre-build
	sls deploy --verbose --stage dev
	

start: build
	sls offline --useDocker start --host 0.0.0.0  --verbose