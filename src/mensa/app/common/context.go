package common

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"strings"

	"github.com/gliderlabs/ssh"
	"github.com/pkg/errors"
)

// 协议类型
type ProtType string

const (
	ProtTypeHTTP ProtType = "http"
	ProtTypeSSH  ProtType = "ssh"
)

type Action string

const (
	ActionTypePull Action = "PULL"
	ActionTypePush Action = "PUSH"
)

func (a Action) IsPull() bool {
	return a == ActionTypePull
}

// 操作者
type Operator struct {
	// 当是http[s]协议时，这里可能有user、pwd(密码可能是token)
	HttpUser *url.Userinfo
	// 当是ssh协议时，可能有ssh的公钥字段
	SSHPublicKey ssh.PublicKey
}

func (o *Operator) IsHttp() bool {
	return o.HttpUser != nil
}

func (o *Operator) IsEmptyUser() bool {
	if o.IsHttp() {
		if o.HttpUser.Username() == "" {
			return true
		}
		if _, set := o.HttpUser.Password(); !set {
			return true
		}
	} else {
		if len(o.SSHPublicKey.Marshal()) == 0 {
			return true
		}
	}
	return false
}

// 相关操作的上下文
type Context struct {
	// push、pull
	ActionType Action
	// 推送方式（http[s]、ssh、git）
	Type ProtType
	// ssh: 原始commands
	RawCommands []string
	// http: 原始url/commands
	RawURL string
	// http: 解析后的url
	RequestURL *url.URL
	// 仓库地址中的owner字段
	RepoOwner string
	// 仓库地址中的 仓库名
	RepoName string
	// 仓库的具体地址
	RepoDir string
	// 推送人 / 拉取人
	// 	当用户提交、拉取仓库时，应该要知道这个操作者是谁
	// 	如果仓库是公共的，那么可以忽略这个操作者字段
	// 	如果仓库是私有的，那么这个字段必须有值
	//
	Operator *Operator

	// http: 请求
	Resp http.ResponseWriter
	Req  *http.Request
}

func (c *Context) Env() []string {
	envSet := make(map[string]string)
	envSet[GROWERLAB_REPO_OWNER] = c.RepoOwner           // 仓库所有者
	envSet[GROWERLAB_REPO_NAME] = c.RepoName             // 仓库名称
	envSet[GROWERLAB_REPO_ACTION] = string(c.ActionType) // 操作类型（pull or push）
	envSet[GROWERLAB_REPO_PROT_TYPE] = string(c.Type)    // 推送方式
	if c.Type == ProtTypeHTTP && c.Operator != nil {
		envSet[GROWERLAB_REPO_OPERATOR] = c.Operator.HttpUser.Username() // 操作者
	} else if c.Type == ProtTypeSSH && c.Operator != nil {
		// TODO 这里之后可能用base64封装一下，目前直接通过grpc到go-git-grpc会有「非法的utf-8错误」
		// envSet[GROWERLAB_REPO_OPERATOR] = string(c.Operator.SSHPublicKey.Marshal())
	}

	result := make([]string, 0, len(envSet))
	for k, v := range envSet {
		result = append(result, fmt.Sprintf("%s=%s", k, v))
	}
	return result
}

func (c *Context) IsReadAction() bool {
	return c.ActionType.IsPull()
}

func (c *Context) Desc() string {
	// who do what
	who := "unknown"
	if c.Operator != nil {
		if c.Operator.IsHttp() {
			who = c.Operator.HttpUser.Username()
		} else {
			who = string(c.Operator.SSHPublicKey.Marshal())
		}
	}
	return fmt.Sprintf("'%s' %s  repo: '%s/%s' on %s", who, c.ActionType, c.RepoOwner, c.RepoName, c.Type)
}

func BuildContextFromHTTP(w http.ResponseWriter, r *http.Request) (*Context, error) {
	uri := r.URL
	repoOwner, repoName, repoPath, err := BuildRepoInfoByPath(uri.Path)
	if err != nil {
		return nil, err
	}

	actionType := ActionTypePush
	service := uri.Query().Get("service")
	if service == "" {
		_, service = path.Split(uri.Path)
	}
	if service == "" {
		return nil, errors.New("invalid service")
	}
	if service == "git-upload-pack" {
		actionType = ActionTypePull
	}

	var operator *Operator = nil
	var username, password, ok = r.BasicAuth()
	if ok {
		operator = &Operator{
			HttpUser: url.UserPassword(username, password),
		}
	}

	return &Context{
		ActionType: actionType,
		Type:       ProtTypeHTTP,
		RawURL:     uri.String(),
		RequestURL: uri,
		RepoOwner:  repoOwner,
		RepoName:   repoName,
		RepoDir:    repoPath, // 仓库的具体地址
		Operator:   operator,
		Resp:       w,
		Req:        r,
	}, nil
}

func BuildContextFromSSH(session ssh.Session) (*Context, error) {
	commands := session.Command()
	if len(commands) < 2 {
		return nil, errors.Errorf("%v commands is invalid", commands)
	}

	gitPath := commands[1]
	repoOwner, repoName, repoPath, err := BuildRepoInfoByPath(gitPath)
	if err != nil {
		return nil, err
	}

	actionType := ActionTypePush
	if commands[0] == "git-upload-pack" {
		actionType = ActionTypePull
	}

	return &Context{
		ActionType:  actionType,
		Type:        ProtTypeSSH,
		RawCommands: commands,
		RepoOwner:   repoOwner,
		RepoName:    repoName,
		RepoDir:     repoPath, // 仓库的地址
		Operator: &Operator{
			SSHPublicKey: session.PublicKey(),
		},
	}, nil
}

func BuildRepoInfoByPath(path string) (repoOwner, repoName, repoPath string, err error) {
	defer func() {
		if e := recover(); e != nil {
			log.Println("build repo info was err: ", e)
		}
	}()

	paths := strings.FieldsFunc(path, func(r rune) bool {
		return r == rune('/') || r == rune('.')
	})
	if len(paths) < 2 {
		err = errors.Errorf("invalid repo path: %s", path)
		return
	}

	repoOwner = paths[0]
	repoName = paths[1]
	repoPath = filepath.Join(repoOwner[:2], repoName[:2], repoOwner, repoName)
	return
}
