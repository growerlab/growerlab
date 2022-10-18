package server

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestObjectStash_Get(t *testing.T) {
	o := NewObjectCache(1 * time.Second)
	o.Set(NewEncodedObject(nil, "123", "", nil))

	obj, ok := o.Get("123")
	assert.NotNil(t, obj)
	assert.Equal(t, true, ok)

	time.Sleep(2 * time.Second)
	obj, ok = o.Get("123")

	assert.Nil(t, obj)
	assert.Equal(t, false, ok)
}
