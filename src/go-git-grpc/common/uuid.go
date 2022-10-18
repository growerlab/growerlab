package common

import (
	"strconv"
	"sync/atomic"
	"time"
)

var id *UUID

func init() {
	id = &UUID{
		base: time.Now().UnixNano(),
	}
	id.start()
}

type UUID struct {
	base     int64
	fakeRand chan int64
}

func (u *UUID) start() {
	u.fakeRand = make(chan int64, 1024)
	go func() {
		for {
			atomic.AddInt64(&u.base, 1)
			u.fakeRand <- u.base
		}
	}()
}

func (u *UUID) Take() string {
	fake := <-u.fakeRand
	id := strconv.FormatInt(fake, 32)
	return id
}

func ShortUUID() string {
	return id.Take()
}
