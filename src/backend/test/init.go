package test

import (
	"os"
	"path"
	"time"

	"github.com/growerlab/growerlab/src/backend/app/model/db"
	namespaceModel "github.com/growerlab/growerlab/src/backend/app/model/namespace"
	repositoryModel "github.com/growerlab/growerlab/src/backend/app/model/repository"
	userModel "github.com/growerlab/growerlab/src/backend/app/model/user"
)

func InitDIR() {
	baseDir := path.Join(os.Getenv("GOPATH"), "src", "github.com/growerlab/growerlab/src/backend")
	err := os.Chdir(baseDir)
	if err != nil {
		panic(err)
	}
}

const (
	FakeEmail    = "molix2@outlook.com"
	FakeUsername = "moli2"

	FakeRepoPath = "hello"
)

func MakeTestRepoData() (*userModel.User, *repositoryModel.Repository, error) {
	user, err := makeUser()
	if err != nil {
		return nil, nil, err
	}

	ns, err := makeNamespace(user)
	if err != nil {
		return nil, nil, err
	}

	repo, err := makeRepository(ns, user)
	if err != nil {
		return nil, nil, err
	}

	return user, repo, err
}

func makeRepository(ns *namespaceModel.Namespace, user *userModel.User) (
	*repositoryModel.Repository, error) {
	repo, _ := repositoryModel.GetRepository(db.DB, 11)
	if repo != nil {
		return repo, nil
	}

	now := time.Now().Unix()
	err := repositoryModel.AddRepository(db.DB, &repositoryModel.Repository{
		ID:          11,
		UUID:        "89F06182F3194B26",
		Path:        FakeRepoPath,
		Name:        FakeRepoPath,
		NamespaceID: ns.ID,
		OwnerID:     user.ID,
		Description: "",
		CreatedAt:   now,
		ServerID:    1,
		ServerPath:  "mo/he/moli2/hello.git",
		Public:      false,
	})
	if err != nil {
		return nil, err
	}
	repo, err = repositoryModel.GetRepository(db.DB, 11)
	return repo, err
}

func makeNamespace(user *userModel.User) (*namespaceModel.Namespace, error) {
	ns, _ := namespaceModel.GetNamespace(db.DB, 2)
	if ns != nil {
		return ns, nil
	}

	err := namespaceModel.AddNamespace(db.DB, &namespaceModel.Namespace{
		ID:      2,
		Path:    "moli2",
		OwnerID: user.ID,
		Type:    int(namespaceModel.TypeUser),
	})
	ns, err = namespaceModel.GetNamespace(db.DB, 2)
	return ns, err
}

func makeUser() (*userModel.User, error) {
	if found, err := userModel.ExistsEmailOrUsername(db.DB, "", FakeEmail); err == nil && found {
		return userModel.GetUserByEmail(db.DB, FakeEmail)
	}

	now := time.Now().Unix()
	err := userModel.AddUser(db.DB, &userModel.User{
		ID:                4,
		Email:             FakeEmail,
		EncryptedPassword: "$argon2id$v=19$m=65536,t=1,p=4$XcRGYkl1YOB5iSy7RqVmzg$E67uTHCFRsT1aMlFg9CTS6QAfhbWAzjjxckFP9JiG+Y",
		Username:          FakeUsername,
		Name:              FakeUsername,
		PublicEmail:       FakeEmail,
		CreatedAt:         now,
		DeletedAt:         nil,
		VerifiedAt:        &now,
		LastLoginAt:       nil,
		LastLoginIP:       nil,
		RegisterIP:        "127.0.0.1",
		IsAdmin:           false,
		NamespaceID:       2,
	})
	if err != nil {
		return nil, err
	}
	currentUser, err := userModel.GetUser(db.DB, 4)
	return currentUser, err
}
