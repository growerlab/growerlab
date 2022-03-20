package app

import (
	"log"

	"github.com/growerlab/growerlab/src/backend/app/common/events"
	"github.com/growerlab/growerlab/src/backend/app/common/notify"
	"github.com/growerlab/growerlab/src/backend/app/common/permission"
	"github.com/growerlab/growerlab/src/backend/app/model/db"
	"github.com/growerlab/growerlab/src/backend/app/utils/conf"
)

// 需要初始化的全局数据放在这里
//	eg. onStart(job.Work)
//
func init() {
	onStart(conf.LoadConfig)
	onStart(db.InitMemDB)
	onStart(db.InitDatabase)
	onStart(notify.InitNotify)
	onStart(permission.InitPermission)
	onStart(events.InitMQ)
}

func onStart(fn func() error) {
	if err := fn(); err != nil {
		log.Printf("%+v\n", err)
		panic(err)
	}
}
