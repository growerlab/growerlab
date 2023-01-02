package repository

import (
	"context"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/growerlab/growerlab/src/backend/app/common/git"
	"github.com/growerlab/growerlab/src/common/errors"
	"github.com/growerlab/growerlab/src/common/path"
)

// TreeFiles 获取仓库的文件列表
// ref: 分支、commit、tag等
// dir: 目录
func (g *Take) TreeFiles(ref string, dir *string) ([]*git.FileEntity, error) {
	if g.repo == nil || govalidator.IsNull(*g.repo) {
		return nil, errors.MissingParameterError(errors.Repository, errors.Repo)
	}
	if dir == nil || govalidator.IsNull(*dir) {
		temp := "/"
		dir = &temp
	}

	var ctx, cancel = context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()

	pathGroup := g.pathGroup()
	files, err := git.New(ctx, pathGroup).TreeFiles(ref, *dir)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return files, nil
}

func (g *Take) pathGroup() string {
	return path.GetPathGroup(g.namespace, *g.repo)
}
