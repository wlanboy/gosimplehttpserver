package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
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

func accesslog(inner http.Handler, logwriter *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		inner.ServeHTTP(w, r)

		logwriter.Printf(
			"%s\t%s",
			r.Method,
			r.RequestURI,
		)
	})
}

func iphandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s", filterPort(getip(r)))
	})
}

func hosthandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s", filterPort(r.Host))
	})
}

func agenthandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s", r.UserAgent())
	})
}

func headershandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for name, value := range r.Header {
			fmt.Fprintf(w, "%v: %v\n", name, value)
		}
	})
}

func main() {

	logger := log.New(os.Stdout, "http: ", log.LstdFlags)

	router := http.NewServeMux()
	router.Handle("/", accesslog(iphandler(), logger))
	router.Handle("/ip", accesslog(iphandler(), logger))
	router.Handle("/agent", accesslog(agenthandler(), logger))
	router.Handle("/host", accesslog(hosthandler(), logger))
	router.Handle("/header", accesslog(headershandler(), logger))

	server := &http.Server{
		Addr:         ":7000",
		ErrorLog:     logger,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
	server.SetKeepAlivesEnabled(false)

	logger.Println("setup signal.Notify")
	closedone := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			logger.Fatalf("Cannot gracefull shutdown %v\n", err)
		}
		close(closedone)
	}()

	logger.Println("ListenAndServe")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Cannot start %v\n", err)
	}
	<-closedone
}
