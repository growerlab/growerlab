package userdomain

import (
	"github.com/growerlab/growerlab/src/backend/app/model/user"
	"github.com/growerlab/growerlab/src/common/errors"
)

type Person struct {
}

func (s *Person) Type() int {
	return TypePerson
}

func (s *Person) TypeLabel() string {
	return "person"
}

func (s *Person) Validate(ud *UserDomain) error {
	if ud.Param <= 0 {
		return errors.Errorf("userdomain param is required")
	}
	return nil
}

func (s *Person) Eval(args Evaluable) ([]int64, error) {
	// TODO 如果这里只是想知道某个用户的 id 的话，那么是可以进行cache的，而不用重复的读取数据库
	u, err := user.GetUser(args.DB().Src, args.UserDomain().Param)
	if err != nil {
		return nil, err
	}
	return []int64{u.ID}, nil
}
