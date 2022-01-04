package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
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

func dumphandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestDump, err := httputil.DumpRequest(r, true)
		if err == nil {
			fmt.Fprint(w, string(requestDump))
		}
	})
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func randomletterString(n int) string {
	var src = rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func pastehandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		code := r.FormValue("code")
		hash := base64.StdEncoding.EncodeToString([]byte(randomletterString(32)))

		file, err := os.Create("./static/" + hash)
		if err == nil {
			file.WriteString(code)
		}
		file.Close()
		if err == nil {
			fmt.Fprintf(w, "<a href='/%v'> your pastebin link </a>\n", hash)
		}
	})
}

func main() {

	logger := log.New(os.Stdout, "http: ", log.LstdFlags)

	router := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./static"))
	router.Handle("/", accesslog(fileServer, logger))

	router.Handle("/ip", accesslog(iphandler(), logger))
	router.Handle("/agent", accesslog(agenthandler(), logger))
	router.Handle("/host", accesslog(hosthandler(), logger))
	router.Handle("/header", accesslog(headershandler(), logger))
	router.Handle("/dump", accesslog(dumphandler(), logger))

	router.Handle("/paste", accesslog(pastehandler(), logger))

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
