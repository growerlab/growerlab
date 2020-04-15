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
	"strings"
)

const (
	DefaultIndex = "index.html"
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

// TODO 后端的接口比较少，可以加个map cache
//
func (w *Router) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	branch := w.branch(req.Host)
	path := DefaultIndex
	if req.URL.Path != "/" {
		path = req.URL.Path
	}
	file := fmt.Sprintf("/data/%s/data/website%s", branch, path)

	log.Printf("url: %s => path: %s client: %s", req.URL.String(), file, req.RemoteAddr)

	if i, err := os.Stat(file); os.IsExist(err) && !i.IsDir() {
		http.ServeFile(resp, req, file)
		return
	}

	// reproxy
	uri, err := url.Parse(fmt.Sprintf("http://services_%s:8081", branch))
	if err != nil {
		panic(errors.New("parse url was err: " + err.Error() + " branch:" + branch))
	}
	reverseProxy := httputil.NewSingleHostReverseProxy(uri)
	// reverseProxy.Director = func(proxyReq *http.Request) {
	// 	// proxyReq.URL.Scheme = uri.Scheme
	// 	// proxyReq.Header = req.Header.Clone()
	// }
	reverseProxy.ModifyResponse = func(response *http.Response) error {
		response.Header.Set("Growerlab", "Router")
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
