package path

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/growerlab/growerlab/src/common/configurator"
	"github.com/growerlab/growerlab/src/common/errors"
)

// GetRealRepositryPath 会根据 pathGroup 获取到真实的仓库地址
func GetRealRepositryPath(pathGroup string) string {
	cfg := configurator.GetConf()
	repoOwner, repoName, repoPath := ParseRepositryPath(pathGroup)
	realPath := filepath.Join(cfg.GitRepoDir, repoOwner, repoName, repoPath)
	return realPath
}

// GetRelativeRepositryPath 会根据 pathGroup 获取到相对仓库地址
func GetRelativeRepositryPath(pathGroup string) string {
	repoOwner, repoName, repoPath := ParseRepositryPath(pathGroup)
	realPath := filepath.Join(repoOwner, repoName, repoPath)
	return realPath
}

// ParseRepositryPath 根据 namespace/username 解析出原始的相对路径
// 加上 config.yaml 中的 git_repo_dir 配置值，即可拿到真实的仓库地址
func ParseRepositryPath(pathGroup string) (repoOwner, repoName, repoPath string) {
	defer func() {
		if e := recover(); e != nil {
			log.Println("build repo info was err: ", e)
		}
	}()

	paths := strings.FieldsFunc(pathGroup, func(r rune) bool {
		return r == rune('/') || r == rune('.')
	})
	if len(paths) < 2 {
		panic(errors.Errorf("invalid repo path: %s", pathGroup))
	}

	repoOwner = paths[0]
	repoName = paths[1]
	repoPath = filepath.Join(repoOwner[:2], repoName[:2], repoOwner, repoName)
	return
}
