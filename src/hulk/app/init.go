package app

import (
	"os"

	"github.com/asaskevich/govalidator"
)

const (
	ActionPull = "pull"
	ActionPush = "push"

	ActionTypeSSH  = "ssh"
	ActionTypeHTTP = "http"
)

var (
	EnvRepoOwner      = os.Getenv("GROWERLAB_REPO_OWNER")
	EnvRepoPath       = os.Getenv("GROWERLAB_REPO_NAME")
	EnvRepoAction     = os.Getenv("GROWERLAB_REPO_ACTION")
	EnvRepoActionType = os.Getenv("GROWERLAB_REPO_PROT_TYPE")
	EnvRepoOperator   = os.Getenv("GROWERLAB_REPO_OPERATOR")
)

func init() {
	switch false {
	case govalidator.IsNull(EnvRepoOwner),
		govalidator.IsNull(EnvRepoPath),
		govalidator.IsNull(EnvRepoAction),
		govalidator.IsNull(EnvRepoActionType),
		govalidator.IsNull(EnvRepoOperator):
		panic("GROWERLAB_REPO_* env variables are not set")
	}
	switch EnvRepoAction {
	case ActionPull, ActionPush:
	default:
		panic("GROWERLAB_REPO_ACTION env variable is not valid")
	}
	switch EnvRepoActionType {
	case ActionTypeSSH, ActionTypeHTTP:
	default:
		panic("GROWERLAB_REPO_PROT_TYPE env variable is not valid")
	}

	ErrPanic(InitConfig())
	ErrPanic(InitRedis())

	app = &App{
		dispatcher: &EventDispatch{},
	}
	app.RegisterHook(&HookEvent{})
}
