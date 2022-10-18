package common

import (
	"fmt"
	"runtime"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUniqueID(t *testing.T) {
	var (
		base    = 1000000
		numCpu  = runtime.NumCPU()
		n       = numCpu * base
		set     = map[string]struct{}{}
		numChan = make(chan string, n)
	)

	assert.LessOrEqual(t, 1, numCpu)

	var wg sync.WaitGroup

	for i := 0; i < numCpu; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < base; j++ {
				id := ShortUUID()
				numChan <- id
			}
		}()
	}

	wg.Wait()

	for {
		select {
		case id := <-numChan:
			if _, exist := set[id]; exist {
				assert.Equal(t, false, exist, fmt.Sprintf("%s was exist", id))
				goto end
			}
			set[id] = struct{}{}
		default:
			goto end
		}
	}

end:
	assert.Equal(t, n, len(set))
}
