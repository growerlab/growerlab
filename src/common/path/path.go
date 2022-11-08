package path

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/growerlab/growerlab/src/common/configurator"
	"github.com/growerlab/growerlab/src/common/errors"
)

// GetRealRepositoryPath 会根据 pathGroup 获取到真实的仓库地址
func GetRealRepositoryPath(pathGroup string) string {
	cfg := configurator.GetConf()
	_, _, repoPath := ParseRepositoryPath(pathGroup)
	realPath := filepath.Join(cfg.GitRepoDir, repoPath)
	return realPath
}

// GetRelativeRepositoryPath 会根据 pathGroup 获取到相对仓库地址
func GetRelativeRepositoryPath(pathGroup string) string {
	_, _, repoPath := ParseRepositoryPath(pathGroup)
	return repoPath
}

// ParseRepositoryPath 根据 namespace/username 解析出原始的相对路径
// 加上 config.yaml 中的 git_repo_dir 配置值，即可拿到真实的仓库地址
func ParseRepositoryPath(pathGroup string) (repoOwner, repoName, repoPath string) {
	defer func() {
		if e := recover(); e != nil {
			log.Println("build repo info was err: ", e)
		}
	}()

	paths := strings.SplitN(pathGroup, "/", 1)
	if len(paths) < 2 {
		panic(errors.Errorf("invalid repo path: %s", pathGroup))
	}

	repoOwner = paths[0]
	repoName = paths[1]
	repoPath = filepath.Join(repoOwner[:2], repoName[:2], repoOwner, repoName)
	return
}

// CheckRepoAbsPathIsEffective 判断某个仓库的路径是否为 config.git_repo_dir 中的路径
// 以此判断 repoPath 的有效性
func CheckRepoAbsPathIsEffective(repoPath string) bool {
	cfg := configurator.GetConf()
	return strings.HasPrefix(repoPath, cfg.GitRepoDir)
}
