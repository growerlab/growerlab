// KeyDB / Redis 配置

package db

import (
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/growerlab/growerlab/src/common/configurator"
	"github.com/growerlab/growerlab/src/common/errors"
)

const KeySep = ":"

var MemDB *MemDBClient

func InitMemDB() error {
	var err error
	var config = configurator.GetConf().Redis

	MemDB, err = DoInitMemDB(config, 0)
	return err
}

func DoInitMemDB(cfg *configurator.Redis, db int) (*MemDBClient, error) {
	mem := newPool(cfg, db)

	// Test
	if err := testMemDB(mem); err != nil {
		return nil, err
	}
	return mem, nil
}

func testMemDB(mem *MemDBClient) error {
	reply, err := mem.Ping().Result()
	if err != nil || reply != "PONG" {
		return errors.New("memdb not ready")
	}
	return nil
}

func newPool(cfg *configurator.Redis, db int) *MemDBClient {
	addr := net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port))
	idleTimeout := time.Duration(cfg.IdleTimeout) * time.Second

	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		DB:           db,
		PoolSize:     cfg.MaxActive,
		MinIdleConns: cfg.MaxIdle,
		IdleTimeout:  idleTimeout,
	})

	memDB := &MemDBClient{
		client,
		NewKeyBuilder(cfg.Namespace),
	}
	return memDB
}

type MemDBClient struct {
	redis.Cmdable
	*KeyBuilder
}

type KeyBuilder struct {
	namespaceKey string
}

func NewKeyBuilder(namespaceKey string) *KeyBuilder {
	return &KeyBuilder{
		namespaceKey: namespaceKey,
	}
}

func (b *KeyBuilder) KeyMaker() *KeyPart {
	var sb = new(strings.Builder)
	if len(b.namespaceKey) > 0 {
		sb.WriteString(b.namespaceKey)
		sb.WriteString(KeySep)
	}

	return &KeyPart{
		sb: sb,
	}
}

// ignore namespace
func (b *KeyBuilder) KeyMakerNoNS() *KeyPart {
	return &KeyPart{
		sb: new(strings.Builder),
	}
}

type KeyPart struct {
	sb *strings.Builder
}

func (b *KeyPart) Append(s ...string) *KeyPart {
	if len(s) == 0 {
		return b
	}

	b.sb.WriteString(s[0])
	for _, k := range s[1:] {
		b.sb.WriteString(KeySep)
		b.sb.WriteString(k)
	}
	return b
}

func (b *KeyPart) String() string {
	return b.sb.String()
}
