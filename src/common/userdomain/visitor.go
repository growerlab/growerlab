package userdomain

type Visitor struct {
}

func (s *Visitor) Type() int {
	return TypeVisitor
}

func (s *Visitor) TypeLabel() string {
	return "everyone"
}

func (s *Visitor) Validate(*UserDomain) error {
	return nil
}

func (s *Visitor) Eval(args Evaluable) ([]int64, error) {
	return []int64{NamespaceVisitor}, nil
}
