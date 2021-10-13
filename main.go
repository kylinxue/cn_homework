package main

import (
	"flag"
	"fmt"
	"os"

	"io"
	"log"
	"net/http"

	"github.com/golang/glog"
)

func main() {
	_ = flag.Set("v", "4")
	flag.Parse()
	glog.Info("Starting http server...")

	http.HandleFunc("/headers", echoHeaders)
	http.HandleFunc("/infos", infos)
	http.HandleFunc("/healthz", healthz)

	err := http.ListenAndServe("0.0.0.0:80", nil)

	if err != nil {
		log.Fatal(err)
	}
}

func infos(w http.ResponseWriter, r *http.Request) {

	env := os.Getenv("JAVA_HOME1")
	host := r.Host
	fmt.Println(host, env)
}

// 返回请求的Headers
func echoHeaders(w http.ResponseWriter, r *http.Request) {
	headers := r.Header
	for k, v := range headers {
		w.Header().Add(k, v[0])
	}
	glog.Info("headers transfer to response!")
}

// 探活
func healthz(w http.ResponseWriter, r *http.Request) {
	_, _ = io.WriteString(w, "ok\n")
}
