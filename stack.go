package main

type stack struct {
	stack []string
}

func (s *stack) Count() int {
	return len(s.stack)
}

func (s *stack) Peep() string {

	l := len(s.stack)

	if l == 0 {
		return ""
	}

	return s.stack[l-1]
}

func (s *stack) Pop() string {

	l := len(s.stack)

	if l == 0 {
		return ""
	}

	v := s.stack[l-1]
	s.stack = s.stack[:l-1]

	return v
}

func (s *stack) Push(v string) {

	s.stack = append(s.stack, v)
}
