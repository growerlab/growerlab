package base

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRange(t *testing.T) {
	data := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11"}

	rg := NewRange(len(data), 2)
	for _, b := range rg.buckets {
		t.Log(b.Start, b.End)
	}

	b1 := rg.buckets[0]
	b2 := rg.buckets[1]
	b3 := rg.buckets[2]
	b4 := rg.buckets[3]
	b5 := rg.buckets[4]
	b6 := rg.buckets[5]

	assert.EqualValues(t, []string{"1", "2"}, data[b1.Start:b1.End])
	assert.EqualValues(t, []string{"3", "4"}, data[b2.Start:b2.End])
	assert.EqualValues(t, []string{"5", "6"}, data[b3.Start:b3.End])
	assert.EqualValues(t, []string{"7", "8"}, data[b4.Start:b4.End])
	assert.EqualValues(t, []string{"9", "10"}, data[b5.Start:b5.End])
	assert.EqualValues(t, []string{"11"}, data[b6.Start:b6.End])

}
