package main

type viewStack struct {
	slice []string
}

func (s *viewStack) push(view string) {
	s.slice = append(s.slice, view)
}

func (s *viewStack) top() string {
	return s.slice[len(s.slice)-1]
}

func (s *viewStack) pop() string {
	top := s.top()
	s.slice = s.slice[:len(s.slice)-1]
	return top
}

func (s *viewStack) replaceTop(view string) {
	s.pop()
	s.push(view)
}

func (s *viewStack) len() int {
	return len(s.slice)
}
