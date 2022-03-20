package reader

import (
	"io"
	"io/ioutil"

	"github.com/growerlab/growerlab/src/backend/app/common/errors"
)

func LimitReader(r io.Reader, n int64) ([]byte, error) {
	lr := io.LimitReader(r, n)
	result, err := ioutil.ReadAll(lr)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return result, nil
}
