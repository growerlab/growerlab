package app

import (
	"github.com/growerlab/growerlab/src/mensa/app/common"
	"github.com/growerlab/growerlab/src/mensa/app/middleware"
)

const (
	DefaultIdleTimeout = 120  // 链接最大闲置时间
	DefaultDeadline    = 3600 // git 的默认执行时间，最长1小时
)

const (
	GitReceivePack   = "git-receive-pack"
	GitUploadPack    = "git-upload-pack"
	GitUploadArchive = "git-upload-archive"

	ReceivePack   = "receive-pack"
	UploadPack    = "upload-pack"
	UploadArchive = "upload-archive"
)

var AllowedCommandMap = map[string]string{
	GitReceivePack:   ReceivePack,
	GitUploadPack:    UploadPack,
	GitUploadArchive: UploadArchive,
}

// 入口
// 	当用户连接到服务
type Entryer interface {
	// 进入前的预备操作
	Enter(ctx *common.Context) (result *middleware.HandleResult)
}
