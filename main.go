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
	Port     string
	BlogURL  string
	PostsURL string

	MicropubEndpoint string
}

var config = &Config{}

func init() {
	config.Port = os.Getenv("PORT")
	if config.Port == "" {
		config.Port = "4567"
	}

	config.BlogURL = os.Getenv("BLOG_URL")
	if config.BlogURL == "" {
		config.BlogURL = "http://localhost:" + config.Port
	}

	config.PostsURL = os.Getenv("POSTS_URL")
	if config.PostsURL == "" {
		config.PostsURL = config.BlogURL + "/posts"
	}

	config.MicropubEndpoint = os.Getenv("MICROPUB_ENDPOINT")
	if config.MicropubEndpoint == "" {
		config.MicropubEndpoint = "https://micro.blog/micropub"
	}
}

func main() {
	router := mux.NewRouter()

	srv := &WPService{config: config}

	codec := xmlrpc.NewCodec()
	codec.AutoCapitalizeMethodName = true

	rs := rpc.NewServer()
	rs.RegisterCodec(codec, "text/xml")
	rs.RegisterService(srv, "wp")
	rs.RegisterService(srv, "metaWeblog")

	router.HandleFunc("/", handleIndex).Methods(http.MethodGet)
	router.HandleFunc("/xmlrpc.php", handleRsd)
	router.Handle("/xmlrpc", rs).Methods(http.MethodPost)

	handler := logHandler(
		handlers.RecoveryHandler()(router),
	)

	addr := ":" + config.Port
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

func handleIndex(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, `<!DOCTYPE html>
<html>
<head>
	<title>sup</title>
	<link rel="EditURI" type="application/rsd+xml" title="RSD" href="%s/xmlrpc.php?rsd" />
</head>
<body>
	<h1>nothing to see here...</h1>
	<p>move along</p>
</body>
</html>`, config.BlogURL)
}

func handleRsd(w http.ResponseWriter, req *http.Request) {
	if req.URL.RawQuery != "rsd" {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	fmt.Fprintf(w, `<rsd xmlns="http://archipelago.phrasewise.com/rsd" version="1.0">
<service>
	<engineName>WordPress</engineName>
	<engineLink>https://wordpress.org/</engineLink>
	<homePageLink>%s</homePageLink>
	<apis>
		<api name="WordPress" blogID="1" preferred="true" apiLink="%s/xmlrpc"/>
	</apis>
	</service>
</rsd>`, config.BlogURL, config.BlogURL)
}
