package jsonutils

import (
	"bytes"
	"encoding/json"

	"github.com/growerlab/growerlab/src/common/errors"
)

func DecodeBytesToObject(payload []byte, out any) error {
	return json.NewDecoder(bytes.NewReader(payload)).Decode(out)
}

func EncodeObjectToBytes(in any) ([]byte, error) {
	w := bytes.NewBuffer(nil)
	err := json.NewEncoder(w).Encode(in)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return w.Bytes(), nil
}
