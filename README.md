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

# calls
* curl http://127.0.0.1:7000/ip
* curl http://127.0.0.1:7000/host
* curl http://127.0.0.1:7000/agent
* curl -H "hello: world" http://127.0.0.1:7000/header
