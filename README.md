![Go](https://github.com/wlanboy/gosimplehttpserver/workflows/Go/badge.svg?branch=master) ![Docker build and publish image](https://github.com/wlanboy/gosimplehttpserver/workflows/Docker%20build%20and%20publish%20image/badge.svg)

# gosimplehttpserver
Golang http server with logging and gracefull shutdown
- 'root' folder with mux and gracefull shutdown
- 'basic' subfolder without mux and gracefull shutdown

# build
* go get -d -v
* go clean
* go build

# run
* go run ip.go

# debug
* go get -u github.com/go-delve/delve/cmd/dlv
* dlv debug ./ip

# go lang build for different boards
- GOOS=linux GOARCH=386 go build (386 needed for busybox)
- GOOS=linux GOARCH=arm GOARM=6 go build (Raspberry Pi build)
- GOOS=linux GOARCH=arm64 go build (Odroid C2 build)

# go lang build for docker
- CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -v .

# docker build
docker build -t gosimplehttpserver:latest . --build-arg BIN_FILE=./gosimplehttpserver

## Docker publish to github registry
- docker tag gosimplehttpserver:latest docker.pkg.github.com/wlanboy/gosimplehttpserver/gosimplehttpserver:latest
- docker push docker.pkg.github.com/wlanboy/gosimplehttpserver/gosimplehttpserver:latest

## Docker hub
- https://hub.docker.com/repository/docker/wlanboy/gosimplehttpserver

## Docker Registry repro
- https://github.com/wlanboy/gosimplehttpserver/packages/278511

# docker run
docker run --name gosimplehttpserver -p 7000:7000 wlanboy/gosimplehttpserver

# calls
* curl http://127.0.0.1:7000/ip
* curl http://127.0.0.1:7000/host
* curl http://127.0.0.1:7000/agent
* curl -H "hello: world" http://127.0.0.1:7000/header
* curl http://127.0.0.1:7000/dump to get the dump of the whole http request object
* http://127.0.0.1:7000/ to get pastebin homepage
* curl -X POST -F 'code=myprettylittleinformation' http://127.0.0.1:7000/paste to create a paste
