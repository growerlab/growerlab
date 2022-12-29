package app

import (
	"net"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/growerlab/growerlab/src/common/db"
)

func InitRedis() error {
	addr := net.JoinHostPort(Conf.Redis.Host, strconv.Itoa(Conf.Redis.Port))
	idleTimeout := time.Duration(Conf.Redis.IdleTimeout) * time.Second

	db.MemDB = &db.MemDBClient{
		UniversalClient: redis.NewClient(&redis.Options{
			Addr:         addr,
			DB:           0,
			PoolSize:     Conf.Redis.MaxActive,
			MinIdleConns: Conf.Redis.MaxIdle,
			IdleTimeout:  idleTimeout,
		}),
		KeyBuilder: db.NewKeyBuilder(Conf.Redis.Namespace),
	}
	return nil
}
