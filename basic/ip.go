package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
)

func filterPort(data string) string {
	if strings.Contains(data, ":") {
		data = strings.Split(data, ":")[0]
	}
	return data
}

func getip(r *http.Request) string {
	ip := r.Header.Get("Http-Client.Ip")
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
	}
	if ip == "" {
		ip = filterPort(r.RemoteAddr)
	}
	return ip
}

func iphandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s\t%s", r.Method, r.RequestURI)
	fmt.Fprintf(w, "%s", filterPort(getip(r)))
}

func hosthandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s\t%s", r.Method, r.RequestURI)
	fmt.Fprintf(w, "%s", filterPort(r.Host))
}

func agenthandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s\t%s", r.Method, r.RequestURI)
	fmt.Fprintf(w, "%s", r.UserAgent())
}

func headershandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s\t%s", r.Method, r.RequestURI)
	for name, value := range r.Header {
		fmt.Fprintf(w, "%v: %v\n", name, value)
	}
}

func dumphandler(w http.ResponseWriter, req *http.Request) {
	requestDump, err := httputil.DumpRequest(req, true)
	if err == nil {
		fmt.Fprint(w, string(requestDump))
	}
}

func main() {
	http.HandleFunc("/", iphandler)
	http.HandleFunc("/ip", iphandler)
	http.HandleFunc("/host", hosthandler)
	http.HandleFunc("/agent", agenthandler)
	http.HandleFunc("/header", headershandler)
	http.HandleFunc("/dump", dumphandler)

	log.Println("start web server")
	log.Fatal(http.ListenAndServe(":7000", nil))
	log.Println("stop web server")
}
