/*
	dev env.
	nginx -> router -> growerlab

	接受nginx的请求，将请求路由到不同的后端 growerlab
*/

package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	DefaultIndex = "index.html"
)

const (
	HeaderID    = "Growerlab"
	HeaderLabel = "Router"
)

var (
	RegexpApi = regexp.MustCompile(`/api/.+`) // api regexp
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	fmt.Println("router running...")
	defer fmt.Println("router exit..")
	return http.ListenAndServe(":81", &Router{})
}

type Router struct {
}

func (w *Router) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	branch := w.branch(req.Host)

	switch {
	case RegexpApi.MatchString(req.URL.Path):
		w.apiLocation(branch, resp, req)
	default:
		w.fileLocation(branch, resp, req)
	}
}

// 文件路由
func (w *Router) fileLocation(branch string, resp http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	root := filepath.Join("/data", branch, "data/website")
	file := filepath.Join(root, path)

	switch path {
	case "/":
		path = DefaultIndex
		file = filepath.Join(root, path)
	default:
		file = filepath.Join(root, path)
		_, err := os.Stat(file)
		if os.IsNotExist(err) {
			path = DefaultIndex
		}
		file = filepath.Join(root, path)
	}

	log.Printf("url: %s, client: %s, path: %s", req.URL.String(), req.RemoteAddr, file)

	http.ServeFile(resp, req, file)
}

// api路由到后端
func (w *Router) apiLocation(branch string, resp http.ResponseWriter, req *http.Request) {
	log.Printf("url: %s, client: %s", req.URL.String(), req.RemoteAddr)

	uri, err := url.Parse(fmt.Sprintf("http://services_%s:8081", branch))
	if err != nil {
		panic(errors.New("parse url was err: " + err.Error() + " branch:" + branch))
	}
	reverseProxy := httputil.NewSingleHostReverseProxy(uri)
	reverseProxy.ModifyResponse = func(response *http.Response) error {
		response.Header.Set(HeaderID, HeaderLabel)
		return nil
	}
	reverseProxy.ServeHTTP(resp, req)
}

func (w *Router) branch(host string) string {
	if !strings.Contains(host, ".dev.growerlab.net") {
		panic(errors.New("invalid host"))
	}
	n := strings.Index(host, ".")
	return host[:n]
}
