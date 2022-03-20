package pwd

import (
	"github.com/growerlab/argon2"
	"github.com/growerlab/growerlab/src/backend/app/common/errors"
)

var argon2Cfg = argon2.DefaultConfig()

func GeneratePassword(src string) (pwd string, err error) {
	raw, err := argon2Cfg.Hash([]byte(src), nil)
	return string(raw.Encode()), errors.Trace(err)
}

func ComparePassword(hashedPwd string, inputPwd string) bool {
	raw, err := argon2.Decode([]byte(hashedPwd))
	if err != nil {
		return false
	}
	b, err := raw.Verify([]byte(inputPwd))
	if err != nil {
		return false
	}
	return b
}
