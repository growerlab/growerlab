package pwd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPassword(t *testing.T) {
	gotPwd, err := GeneratePassword("hello pwd")
	if err != nil {
		t.Fatal(err)
	}

	ok := ComparePassword(gotPwd, "hello pwd")
	assert.Equal(t, true, ok, nil)
}
