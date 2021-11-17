package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"io"
	"log"
	"net/http"
)

func main() {
	//_ = flag.Set("v", "4")
	//flag.Parse()
	log.Println("Starting http server...")

	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/headers", echoHeaders)
	mux.HandleFunc("/infos", infos)
	mux.HandleFunc("/healthz", healthz)

	srv := http.Server{
		Addr:    ":80",
		Handler: mux,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Println("Server Started")

	<-done
	log.Println("Server Stopping...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")

}

func infos(w http.ResponseWriter, r *http.Request) {

	env := os.Getenv("JAVA_HOME1")
	host := r.Host
	fmt.Println(host, env)
}

// 返回请求的Headers
func echoHeaders(w http.ResponseWriter, r *http.Request) {
	headers := r.Header
	for k, vs := range headers {
		for _, v := range vs {
			w.Header().Add(k, v)
		}
	}
	log.Println("headers transfer to response!")
}

// 探活
func healthz(w http.ResponseWriter, r *http.Request) {
	_, _ = io.WriteString(w, "ok\n")
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("entering root handler")
	user := r.URL.Query().Get("user")
	if user != "" {
		io.WriteString(w, fmt.Sprintf("hello [%s]\n", user))
	} else {
		io.WriteString(w, "hello [stranger]\n")
	}
	io.WriteString(w, "===================Details of the http request header:============\n")
	for k, v := range r.Header {
		io.WriteString(w, fmt.Sprintf("%s=%s\n", k, v))
	}
}
