package common

type Enum struct {
	labelSet map[int]string
	enumSet  map[string]int
}

func NewEnum() *Enum {
	return &Enum{
		labelSet: map[int]string{},
		enumSet:  map[string]int{},
	}
}

func (e *Enum) Set(id int, label string) {
	e.labelSet[id] = label
	e.enumSet[label] = id
}

func (e *Enum) Label(id int) (string, bool) {
	s, ok := e.labelSet[id]
	return s, ok
}

func (e *Enum) ID(label string) (int, bool) {
	n, ok := e.enumSet[label]
	return n, ok
}
