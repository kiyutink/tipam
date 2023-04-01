package main

func newStack[T any]() *stack[T] {
	return &stack[T]{
		slice: []T{},
	}
}

type stack[T any] struct {
	slice []T
}

func (s *stack[T]) push(t T) {
	s.slice = append(s.slice, t)
}

func (s *stack[T]) top() T {
	return s.slice[len(s.slice)-1]
}

func (s *stack[T]) pop() T {
	top := s.top()
	s.slice = s.slice[:len(s.slice)-1]
	return top
}

func (s *stack[T]) replaceTop(t T) {
	s.pop()
	s.push(t)
}

func (s *stack[T]) len() int {
	return len(s.slice)
}
