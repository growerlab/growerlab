package service

import (
	"fmt"
	"strconv"

	repoModel "github.com/growerlab/growerlab/src/backend/app/model/repository"
	"github.com/growerlab/growerlab/src/common/db"
	"github.com/growerlab/growerlab/src/common/errors"
)

var NotFoundRepoError = errors.New("not found repo")

func GetRepository(repoOwner, repoName string) (*repoModel.Repository, error) {
	repoOwnerNS, err := GetUserNamespaceByUsername(repoOwner)
	if err != nil {
		return nil, err
	}

	key := db.MemDB.KeyMaker().Append("repository", "id", "namespace").String()
	field := db.MemDB.KeyMakerNoNS().Append(repoOwner, repoName).String()

	// 仓库的公开状态等属性可能变动，所以这里仅缓存仓库id
	repoIDRaw, err := NewCache().GetOrSet(
		key,
		field,
		func() (value string, err error) {
			repo, err := repoModel.GetRepositoryByNsWithPath(db.DB, repoOwnerNS, repoName)
			if err != nil {
				return "", err
			}
			if repo == nil {
				return "", errors.Message(NotFoundRepoError, fmt.Sprintf("%s/%s", repoOwner, repoName))
			}
			return strconv.FormatInt(repo.ID, 10), nil
		})
	if err != nil {
		return nil, err
	}

	repoID, err := strconv.ParseInt(repoIDRaw, 10, 64)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return repoModel.GetRepository(db.DB, repoID)
}
