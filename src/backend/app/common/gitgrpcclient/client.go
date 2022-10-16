package gitgrpcclient

import (
	"context"

	"github.com/growerlab/growerlab/src/common/errors"

	"github.com/growerlab/growerlab/src/backend/app/common/notify"
	"github.com/growerlab/growerlab/src/common/configurator"
	gggrpc "github.com/growerlab/growerlab/src/go-git-grpc"
	"github.com/growerlab/growerlab/src/go-git-grpc/client"
)

func GetGitGRPCClient(ctx context.Context, repoPath string) (*client.Store, error) {
	global := configurator.GetConf()
	store, closeFn, err := gggrpc.NewStoreClient(ctx, global.GoGitGrpcServerAddr, repoPath)
	if err != nil {
		return nil, errors.Trace(err)
	}
	notify.Subscribe(func() {
		closeFn()
	})

	return store, nil
}
