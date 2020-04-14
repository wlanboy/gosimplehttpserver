![Go](https://github.com/wlanboy/gosimplehttpserver/workflows/Go/badge.svg?branch=master)

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

# go lang build for docker
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -v .

# docker build
docker build -t gosimplehttpserver:latest . --build-arg BIN_FILE=./gosimplehttpserver

# docker run
docker run --name kanban -p 7000:7000 wlanboy/gosimplehttpserver

# calls
* curl http://127.0.0.1:7000/ip
* curl http://127.0.0.1:7000/host
* curl http://127.0.0.1:7000/agent
* curl -H "hello: world" http://127.0.0.1:7000/header
