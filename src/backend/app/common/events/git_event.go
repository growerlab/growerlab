package events

import (
	"fmt"
	"github.com/growerlab/growerlab/src/common/errors"

	"github.com/growerlab/growerlab/src/backend/app/common/mq"
)

type PushSession struct {
	RefName string `json:"refname"` // branch, tag name
	OldRev  string `json:"oldrev"`  // old
	NewRev  string `json:"newrev"`  // new

	RepoDir string `json:"repo_dir"`

	RepoOwner string `json:"repo_owner"` // namespace.path
	RepoPath  string `json:"repo_path"`  // repository name

	Action     string `json:"action"`
	ActionType string `json:"action_type"` // 当tag提交到服务器时有/无描述的特征
	RefType    string `json:"ref_type"`

	ProtType string `json:"prot_type"` // http/ssh

	Opeator string `json:"operator"` // 推送者
}

type Signature struct {
	// Name represents a person name. It is an arbitrary string.
	Name string `json:"name"`
	// Email is an email, but it cannot be assumed to be well-formed.
	Email string `json:"email"`
	// When is the timestamp of the signature.
	When int64 `json:"when"`
}

type PlainCommit struct {
	Hash      string    `json:"hash"`
	Author    Signature `json:"author"`
	Committer Signature `json:"committer"`
	Message   string    `json:"message"`
}

type GitEventPayload struct {
	Session     *PushSession   `json:"session"`
	CommitCount int            `json:"commit_count"`
	Commits     []*PlainCommit `json:"commits"`        // commits
	Message     string         `json:"commit_message"` // commit/tag message
}

type AsyncPushGitEvent interface {
	AsyncPushGitEvent(gitEvent interface{}) error
}

var _ mq.Consumer = (*GitEvent)(nil)

type GitEvent struct {
}

func newGitEventConsumer() mq.Consumer {
	return &GitEvent{}
}

func NewGitEvent() AsyncPushGitEvent {
	return &GitEvent{}
}

func (*GitEvent) Name() string {
	return "git_event"
}

func (*GitEvent) DefaultField() string {
	return "default"
}

func (g *GitEvent) Consume(payload *mq.Payload) error {
	gitEventPayload := new(GitEventPayload)
	err := getPayload(payload, g.DefaultField(), gitEventPayload)
	if err != nil {
		return errors.Trace(err)
	}
	// TODO 消费消息
	fmt.Println(gitEventPayload)
	return nil
}

func (g *GitEvent) AsyncPushGitEvent(gitEvent interface{}) error {
	return async(g.Name(), g.DefaultField(), gitEvent)
}
