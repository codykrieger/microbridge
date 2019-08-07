package main

import (
	// "bytes"
	"fmt"
	// "io/ioutil"
	"net/http"
	"os"
	"time"

	xmlrpc "github.com/codykrieger/gorilla-xmlrpc/xml"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	log "github.com/sirupsen/logrus"
)

func main() {
	router := mux.NewRouter()

	srv := &WPService{}

	codec := xmlrpc.NewCodec()
	codec.AutoCapitalizeMethodName = true

	rs := rpc.NewServer()
	rs.RegisterCodec(codec, "text/xml")
	rs.RegisterService(srv, "wp")

	router.Handle("/xmlrpc", rs).Methods(http.MethodPost)

	port := os.Getenv("PORT")
	if port == "" {
		port = "4567"
	}

	addr := ":" + port

	handler := logHandler(
		handlers.RecoveryHandler()(router),
	)

	log.WithField("addr", addr).Info("listening...")

	if err := http.ListenAndServe(addr, handler); err != nil {
		fatalf("http.ListenAndServe: %v", err)
	}
}

func fatalf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "fatal: "+format+"\n", args...)
	os.Exit(1)
}

type loggingResponseWriter struct {
	w      http.ResponseWriter
	status int
	size   int
}

func (l *loggingResponseWriter) Header() http.Header {
	return l.w.Header()
}

func (l *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := l.w.Write(b)
	l.size += size
	return size, err
}

func (l *loggingResponseWriter) WriteHeader(statusCode int) {
	l.w.WriteHeader(statusCode)
	l.status = statusCode
}

func logHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		t := time.Now()

		// if buf, err := ioutil.ReadAll(req.Body); err != nil {
		//     http.Error(w, err.Error(), http.StatusInternalServerError)
		//     return
		// } else {
		//     fmt.Fprintf(os.Stderr, "\n\033[33mbody: %s\033[0m\n\n", ioutil.NopCloser(bytes.NewBuffer(buf)))
		//     req.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
		// }

		lrw := &loggingResponseWriter{w: w, status: http.StatusOK}
		next.ServeHTTP(lrw, req)

		d := time.Since(t)

		log.WithFields(log.Fields{
			"code":   lrw.status,
			"sz":     req.ContentLength,
			"respsz": lrw.size,
			"d":      d,
			"ct":     req.Header.Get("content-type"),
			"ua":     req.Header.Get("user-agent"),
		}).Info(req.Method + " " + req.RequestURI)
	})
}
