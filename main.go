package main

import (
	// "bytes"
	"fmt"
	// "io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/codykrieger/microbridge/xmlrpc"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	rpc "github.com/gorilla/rpc/v2"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	BlogURL  string
	PostsURL string

	Endpoint string
	Username string
	// Token    string
}

var config = &Config{}

func init() {
	config.BlogURL = os.Getenv("BLOG_URL")
	if config.BlogURL == "" {
		// panic("BLOG_URL environment variable must be set")
		config.BlogURL = "https://cjk.micro.blog"
	}

	config.PostsURL = os.Getenv("POSTS_URL")
	if config.PostsURL == "" {
		config.PostsURL = config.BlogURL + "/posts"
	}

	config.Endpoint = os.Getenv("API_ENDPOINT")
	if config.Endpoint == "" {
		config.Endpoint = "https://micro.blog/micropub"
	}

	config.Username = os.Getenv("API_USER")
	if config.Username == "" {
		config.Username = "you"
	}

	// config.Token = os.Getenv("API_TOKEN")
	// if config.Token == "" {
	//     // panic("API_TOKEN environment variable must be set")
	//     config.Token = "XXXXXXXXXXXXXXXXXXXX"
	// }
}

func main() {
	router := mux.NewRouter()

	srv := &WPService{config: config}

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

		msg := req.Method + " " + req.RequestURI
		l := log.WithFields(log.Fields{
			"code":   lrw.status,
			"sz":     req.ContentLength,
			"respsz": lrw.size,
			"d":      d,
			"ct":     req.Header.Get("content-type"),
			"ua":     req.Header.Get("user-agent"),
		})

		if lrw.status >= 200 && lrw.status < 300 {
			l.Info(msg)
		} else if lrw.status >= 300 && lrw.status < 400 {
			l.Warn(msg)
		} else {
			l.Error(msg)
		}
	})
}
