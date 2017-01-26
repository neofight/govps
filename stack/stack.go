package stack

type Stack struct {
	stack []string
}

func (s *Stack) Count() int {
	return len(s.stack)
}

func (s *Stack) Peep() string {

	l := len(s.stack)

	if l == 0 {
		return ""
	}

	return s.stack[l-1]
}

func (s *Stack) Pop() string {

	l := len(s.stack)

	if l == 0 {
		return ""
	}

	v := s.stack[l-1]
	s.stack = s.stack[:l-1]

	return v
}

func (s *Stack) Push(v string) {

	s.stack = append(s.stack, v)
}
