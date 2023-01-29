package repository

import (
	"context"
	"path/filepath"
	"sort"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/growerlab/growerlab/src/backend/app/common/git"
	"github.com/growerlab/growerlab/src/common/errors"
	"github.com/growerlab/growerlab/src/common/path"
)

// TreeFiles 获取仓库的文件列表
// ref: 分支、commit、tag等
// dir: 目录
func (g *Take) TreeFiles(ref string, folder *string) ([]*git.FileEntity, error) {
	if g.repo == nil || govalidator.IsNull(*g.repo) {
		return nil, errors.MissingParameterError(errors.Repository, errors.Repo)
	}
	if folder == nil || govalidator.IsNull(*folder) {
		temp := "/"
		folder = &temp
	}
	*folder = filepath.Clean(*folder)

	var ctx, cancel = context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()

	pathGroup := g.pathGroup()
	files, err := git.New(ctx, pathGroup).TreeFiles(ref, *folder)
	if err != nil {
		if err.Error() == object.ErrDirectoryNotFound.Error() {
			return nil, errors.NotFoundError(errors.Folder)
		}
		return nil, errors.Trace(err)
	}

	// sort
	sort.Sort(git.FileEntitySorter(files))
	return files, nil
}

func (g *Take) pathGroup() string {
	return path.GetPathGroup(g.namespace, *g.repo)
}
