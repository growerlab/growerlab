/*
	dev env.
	nginx -> router -> growerlab

	接受nginx的请求，将请求路由到不同的后端 growerlab
*/

package app

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
	"time"
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

func RunWebSiteRouter() {
	go func() {
		if err := runWebSiteRouter(); err != nil {
			panic(err)
		}
	}()
}

func runWebSiteRouter() error {
	fmt.Println("router running...")
	defer fmt.Println("router exit..")
	return http.ListenAndServe(":8010", &WebSiteRouter{})
}

type WebSiteRouter struct {
}

func (w *WebSiteRouter) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	branch := BranchByHost(req.Host)

	switch {
	case RegexpApi.MatchString(req.URL.Path):
		w.apiLocation(branch, resp, req)
	default:
		w.fileLocation(branch, resp, req)
	}
}

// 文件路由
func (w *WebSiteRouter) fileLocation(branch string, resp http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	root := filepath.Join("/data", branch, "data/website")
	file := filepath.Join(root, path)

	modtime := time.Now()
	var fp *os.File

	switch path {
	case "/":
		path = DefaultIndex
		file = filepath.Join(root, path)
		stat, err := os.Stat(file)
		if os.IsNotExist(err) {
			resp.WriteHeader(http.StatusNotFound)
			return
		}
		modtime = stat.ModTime()
		fp, err = os.Open(file)
		if err != nil {
			resp.Write([]byte(err.Error()))
			return
		}
	default:
		file = filepath.Join(root, path)
		stat, err := os.Stat(file)
		if os.IsNotExist(err) {
			path = DefaultIndex
		}
		file = filepath.Join(root, path)
		modtime = stat.ModTime()
		fp, err = os.Open(file)
		if err != nil {
			resp.Write([]byte(err.Error()))
			return
		}
	}

	log.Printf("url: %s, client: %s, path: %s", req.URL.String(), req.RemoteAddr, file)

	http.ServeContent(resp, req, file, modtime, fp)
}

// api路由到后端
func (w *WebSiteRouter) apiLocation(branch string, resp http.ResponseWriter, req *http.Request) {
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
