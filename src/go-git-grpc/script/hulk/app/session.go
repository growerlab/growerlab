package app

import (
	"encoding/json"
	repo2 "github.com/growerlab/growerlab/src/go-git-grpc/script/hulk/app/repo"
	"os"

	"github.com/asaskevich/govalidator"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/growerlab/growerlab/src/common/errors"
)

type Action string
type ActionType string
type RefType string

const (
	ActionCreate Action = "create" // create branch or tag
	ActionDelete Action = "delete" // delete branch or tag
	ActionCommit Action = "commit" // push commit
)

const (
	ActionTypeUnannotated = "unannotated" // 没有描述的tag
	ActionTypeAnnotated   = "annotated"   // 有描述的tag
)

const (
	RefTypeBranch RefType = "branch"
	RefTypeTag    RefType = "tag"
)

type PushSession struct {
	RefName string `json:"refname"` // branch, tag name
	OldRev  string `json:"oldrev"`  // old
	NewRev  string `json:"newrev"`  // new

	RepoDir string `json:"repo_dir"`

	RepoOwner string `json:"repo_owner"` // namespace.path
	RepoPath  string `json:"repo_path"`  // repository name

	Action     Action  `json:"action"`
	ActionType string  `json:"action_type"` // 当tag提交到服务器时有/无描述的特征
	RefType    RefType `json:"ref_type"`

	ProtType string `json:"prot_type"` // http/ssh

	Opeator string `json:"operator"` // 推送者
}

func (r *PushSession) JSON() string {
	if r == nil {
		return ""
	}
	b, _ := json.Marshal(r)
	return string(b)
}

func (r *PushSession) IsNullOldCommit() bool {
	return r.OldRev == repo2.ZeroRef
}

func (r *PushSession) IsNullNewCommit() bool {
	return r.NewRev == repo2.ZeroRef
}

func (r *PushSession) IsNewBranch() bool {
	return r.RefType == RefTypeBranch && r.Action == ActionCreate
}

func (r *PushSession) IsNewTag() bool {
	return r.RefType == RefTypeTag && r.Action == ActionCreate
}

func (r *PushSession) IsDeleteAction() bool {
	return r.IsNullNewCommit()
}

func (r *PushSession) IsCommitPush() bool {
	return !r.IsNullOldCommit() && !r.IsNullNewCommit()
}

func (r *PushSession) RevType(rev string) string {
	if rev == repo2.ZeroRef {
		return ""
	}
	t := repo2.NewRepository(r.RepoDir).HashType(plumbing.NewHash(rev))
	switch t {
	case plumbing.TagObject:
		return "tag"
	case plumbing.CommitObject:
		return "commit"
	}
	return ""
}

func (r *PushSession) prepare() error {
	if govalidator.IsNull(r.RefName) {
		return errors.New("ref name is empty")
	}
	if govalidator.IsNull(r.OldRev) {
		return errors.New("old rev is empty")
	}
	if govalidator.IsNull(r.NewRev) {
		return errors.New("new rev is empty")
	}

	r.RepoOwner = EnvRepoOwner
	r.RepoPath = EnvRepoPath

	if repo2.IsTag(r.RefName) {
		r.RefType = RefTypeTag
		r.Action = ActionCreate
		if r.IsDeleteAction() {
			r.Action = ActionDelete
		} else {
			r.ActionType = ActionTypeAnnotated
			if r.RevType(r.NewRev) == "commit" {
				r.ActionType = ActionTypeUnannotated
			}
		}
	} else if repo2.IsBranch(r.RefName) {
		r.RefType = RefTypeBranch
		if r.IsDeleteAction() {
			r.Action = ActionDelete
		} else {
			if r.RevType(r.NewRev) == "commit" {
				r.Action = ActionCreate
			}
		}
	} else if r.IsCommitPush() {
		r.RefType = RefTypeBranch
		r.Action = ActionCommit
	} else {
		return errors.Errorf("invalid ref '%s'", r.RefName)
	}
	return nil
}

func Session() *PushSession {
	pwd, err := os.Getwd()
	if err != nil {
		ErrPanic(err)
	}

	sess := &PushSession{
		RepoDir: pwd,
		RefName: os.Args[1],
		OldRev:  os.Args[2],
		NewRev:  os.Args[3],
	}

	if err := sess.prepare(); err != nil {
		ErrPanic(err)
	}
	return sess
}

// ErrPanic non-zero exit code
func ErrPanic(err error) {
	if err != nil {
		panic(errors.WithStack(err))
	}
}
