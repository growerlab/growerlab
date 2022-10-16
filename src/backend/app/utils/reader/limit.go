package reader

import (
	"github.com/growerlab/growerlab/src/common/errors"
	"io"
	"io/ioutil"
)

func LimitReader(r io.Reader, n int64) ([]byte, error) {
	lr := io.LimitReader(r, n)
	result, err := ioutil.ReadAll(lr)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return result, nil
}
