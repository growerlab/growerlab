package server

import (
	"log"
	"time"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/growerlab/growerlab/src/go-git-grpc/common"
	gocache "github.com/patrickmn/go-cache"
)

const (
	defaultObjectTimeOut   = 30 * time.Second
	defaultCleanupInterval = 10 * time.Minute
)

type Object interface {
	UUID() string
}

type ObjectCache struct {
	cache *gocache.Cache
}

func NewObjectCache(objectTimeout time.Duration) *ObjectCache {
	if objectTimeout <= 0 {
		objectTimeout = defaultObjectTimeOut
	}

	return &ObjectCache{
		cache: gocache.New(objectTimeout, defaultCleanupInterval),
	}
}

func (c *ObjectCache) Set(obj Object) {
	key := obj.UUID()
	c.cache.SetDefault(key, obj)
}

func (c *ObjectCache) Get(uuid string) (Object, bool) {
	ee, ok := c.cache.Get(uuid)
	if !ok {
		return nil, false
	}
	obj, ok := ee.(Object)
	if !ok {
		return nil, false
	}
	c.Set(obj) // 延期
	return obj, true
}

func buildUUID(obj plumbing.EncodedObject) string {
	if obj == nil {
		return common.ShortUUID()
	}
	switch obj.Type() {
	case plumbing.CommitObject,
		plumbing.TreeObject,
		plumbing.BlobObject,
		plumbing.TagObject:
		return obj.Hash().String()
	case plumbing.OFSDeltaObject:
		log.Panic("==== 注意：这里暂时不知返回什么uuid")
	case plumbing.REFDeltaObject:
		log.Panic("==== 注意：这里暂时不知返回什么uuid")
	}
	return common.ShortUUID()
}
