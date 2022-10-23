package git

import (
	"context"
	"io"

	"github.com/growerlab/growerlab/src/common/configurator"
	"github.com/growerlab/growerlab/src/common/errors"
	gggrpc "github.com/growerlab/growerlab/src/go-git-grpc"
	"github.com/growerlab/growerlab/src/go-git-grpc/client"
)

func GetGitGRPCClient(ctx context.Context, repoRelativePath string) (*client.Store, io.Closer, error) {
	global := configurator.GetConf()
	store, closeFn, err := gggrpc.NewStoreClient(ctx, global.GoGitGrpcServerAddr, repoRelativePath)
	if err != nil {
		return nil, nil, errors.Trace(err)
	}
	return store, closeFn, nil
}
