package context

import (
	"github.com/growerlab/growerlab/src/backend/app/common/errors"
)

type Repository struct {
}

func (r *Repository) Type() int {
	return TypeRepository
}

func (r *Repository) TypeLabel() string {
	return "repository"
}

func (r *Repository) Validate(c *Context) error {
	if c.Param1 <= 0 {
		return errors.Errorf("context param1 is required")
	}
	return nil
}
