package app

import (
	"compress/gzip"
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/growerlab/growerlab/src/common/configurator"
	"github.com/growerlab/growerlab/src/common/errors"
	"github.com/growerlab/growerlab/src/common/path"
	"github.com/growerlab/growerlab/src/mensa/app/common"
)

const BannerMessage = "----- Power by GrowerLab.net -----"

func NewGitHttpServer(cfg *configurator.Mensa) *GitHttpServer {
	deadline := DefaultDeadline * time.Second
	idleTimeout := DefaultIdleTimeout * time.Second

	if cfg.Deadline > 0 {
		deadline = time.Duration(cfg.Deadline) * time.Second
	}
	if cfg.IdleTimeout > 0 {
		idleTimeout = time.Duration(cfg.IdleTimeout) * time.Second
	}

	server := &GitHttpServer{
		listen:      cfg.HTTPListen,
		gitBinPath:  cfg.GitPath,
		deadline:    deadline,
		idleTimeout: idleTimeout,
	}

	engine := gin.Default()
	engine.Use(server.handlerBuildRequestContext)
	engine.GET("/:path/:repo_name/info/refs", server.handlerGetInfoRefs)
	engine.POST("/:path/:repo_name/:rpc", server.handlerGitRPC)

	server.server = &http.Server{
		Handler:      engine,
		Addr:         cfg.HTTPListen,
		WriteTimeout: deadline,
		ReadTimeout:  deadline,
		IdleTimeout:  idleTimeout,
	}
	return server
}

type requestContext struct {
	c        *gin.Context
	w        http.ResponseWriter
	r        *http.Request
	Rpc      string
	RepoPath string
}

type GitHttpServer struct {
	// engine for http git
	server *http.Server
	// 服务器的监听地址(eg. host:port)
	listen string
	// git bin path
	gitBinPath string
	// logger io.Writer
	// 最长执行时间
	deadline time.Duration
	// 限制最大时间
	idleTimeout time.Duration

	middlewareHandler MiddlewareHandler
}

func (g *GitHttpServer) ListenAndServe(handler MiddlewareHandler) error {
	log.Printf("[http] git listen and serve: %v\n", g.listen)

	if err := g.validate(); err != nil {
		return err
	}

	g.middlewareHandler = handler

	if err := g.server.ListenAndServe(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (g *GitHttpServer) handlerBuildRequestContext(c *gin.Context) {
	r := c.Request
	w := c.Writer
	// file := r.URL.Path
	_, _, repoPath := path.ParseRepositryPath(r.URL.Path)

	rpc := g.getServiceType(c)

	req := &requestContext{
		c:        c,
		w:        w,
		r:        r,
		Rpc:      rpc,
		RepoPath: repoPath,
	}
	c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "request_context", req))
}

func (g *GitHttpServer) handlerGitRPC(c *gin.Context) {
	reqContext, ok := c.Request.Context().Value("request_context").(*requestContext)
	if !ok {
		log.Println("handlerGitRPC: 'request_context' must exist in context")
		c.AbortWithStatus(http.StatusServiceUnavailable)
		return
	}

	err := g.serviceRpc(reqContext)
	if err != nil {
		log.Printf("git rpc err: %+v\n", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}

func (g *GitHttpServer) handlerGetInfoRefs(c *gin.Context) {
	reqContext, ok := c.Request.Context().Value("request_context").(*requestContext)
	if !ok {
		log.Println("handlerGetInfoRefs: 'request_context' must exist in context")
		c.AbortWithStatus(http.StatusServiceUnavailable)
		return
	}

	err := g.getInfoRefs(reqContext)
	if err != nil {
		log.Printf("get info refs was err: %v\n", err)
		return
	}
}

func (g *GitHttpServer) validate() error {
	if g.listen == "" {
		return errors.New("addr is required")
	}
	if !strings.Contains(g.listen, ":") {
		return errors.Errorf("addr is invalid: %s", g.listen)
	}
	return nil
}

// 平滑重启
func (g *GitHttpServer) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return g.server.Shutdown(ctx)
}

func (g *GitHttpServer) runMiddlewares(ctx *common.Context) error {
	result := g.middlewareHandler(ctx)
	if result.HttpCode > http.StatusCreated {
		ctx.Resp.WriteHeader(result.HttpCode)
		ctx.Resp.Header().Set("WWW-Authenticate", "Basic") // fmt.Sprintf("Basic realm=%s charset=UTF-8"))
	}

	if result.Err != nil {
		log.Printf("[http] middleware err: %+v \nresult:%d %s\n", result.Err, result.HttpCode, result.HttpMessage)
	}
	return result.Err
}

func (g *GitHttpServer) serviceRpc(ctx *requestContext) error {
	var w, r, rpc, dir = ctx.w, ctx.r, ctx.Rpc, ctx.RepoPath

	var body = r.Body
	defer body.Close()

	commonCtx, err := common.BuildContextFromHTTP(ctx.w, ctx.r)
	if err != nil {
		return errors.WithStack(err)
	}

	w.Header().Set("Content-Type", fmt.Sprintf("application/x-git-%s-result", rpc))

	if r.Header.Get("Content-Encoding") == "gzip" {
		body, _ = gzip.NewReader(r.Body)
	}

	// 客户端push：输出到客户端的终端，之后这块应该要抽出来结构化
	if rpc == ReceivePack {
		_, _ = w.Write(packetWrite(fmt.Sprintf("\x02 %s\n", BannerMessage)))
	}

	args := make([]string, 0)
	if rpc == ReceivePack {
		for _, op := range GitReceivePackOptions {
			args = append(args, op.Name, op.Args)
		}
	}
	args = append(args, rpc, "--stateless-rpc", ".")

	err = gitCommand(body, w, dir, args, commonCtx.Env())
	if err != nil {
		log.Printf("git command err: %+v\n", err)
		return err
	}

	// 当有修改仓库时，更新仓库
	if rpc == ReceivePack {
		err = updateServerInfo(dir, commonCtx.Env())
	}

	return errors.WithStack(err)
}

func (g *GitHttpServer) getInfoRefs(ctx *requestContext) error {
	w, r, rpc, dir := ctx.w, ctx.r, ctx.Rpc, ctx.RepoPath

	access := g.hasAccess(r, dir, rpc, false)
	if access {
		commonCtx, err := common.BuildContextFromHTTP(ctx.w, ctx.r)
		if err != nil {
			return errors.WithStack(err)
		}
		err = g.runMiddlewares(commonCtx)
		if err != nil {
			return err
		}

		g.hdrNocache(w)
		w.Header().Set("Content-Type", fmt.Sprintf("application/x-git-%s-advertisement", rpc))
		_, _ = w.Write(packetWrite("# service=git-" + rpc + "\n"))
		_, _ = w.Write(packetFlush())

		args := []string{rpc, "--stateless-rpc", "--advertise-refs", "."}
		err = gitCommand(r.Body, w, dir, args, commonCtx.Env())
		if err != nil {
			return err
		}
	} else {
		log.Printf("can't access %s %s\n", dir, rpc)
	}
	return nil
}

func packetFlush() []byte {
	return []byte("0000")
}

func packetWrite(str string) []byte {
	s := strconv.FormatInt(int64(len(str)+4), 16)

	if len(s)%4 != 0 {
		s = strings.Repeat("0", 4-len(s)%4) + s
	}
	return []byte(s + str)
}

// hasAccess 是否有访问权限
// 这里之后可能要改成从数据库、权限验证中心来确认
func (g *GitHttpServer) hasAccess(r *http.Request, dir string, rpc string, checkContentType bool) bool {
	if checkContentType {
		if r.Header.Get("Content-Type") != fmt.Sprintf("application/x-git-%s-request", rpc) {
			return false
		}
	}

	if !(rpc == UploadPack || rpc == ReceivePack) {
		return false
	}
	if rpc == ReceivePack {
		// return g.config.ReceivePack
		return true
	}
	if rpc == UploadPack {
		// return g.config.UploadPack
		return true
	}

	return true
}

func (g *GitHttpServer) getServiceType(c *gin.Context) string {
	serviceType := c.Request.FormValue("service")
	if len(serviceType) == 0 {
		serviceType = c.Param("rpc")
	}

	if s := strings.HasPrefix(serviceType, "git-"); !s {
		return ""
	}
	return strings.Replace(serviceType, "git-", "", 1)
}

func (g *GitHttpServer) getGitConfig(configName string, dir string) (string, error) {
	var args = []string{"config", configName}
	var out strings.Builder
	err := gitCommand(nil, &out, dir, args, nil)
	return out.String(), errors.WithStack(err)
}

func (g *GitHttpServer) hdrNocache(w http.ResponseWriter) {
	w.Header().Set("Expires", "Fri, 01 Jan 1980 00:00:00 GMT")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Cache-Control", "no-cache, max-age=0, must-revalidate")
}

func updateServerInfo(dir string, envs []string) error {
	err := gitCommand(nil, nil, dir, []string{"--git-dir", ".", "update-server-info"}, envs)
	if err != nil {
		log.Printf("git command 'update-server-info' err: %+v\n", err)
	}
	return err
}
