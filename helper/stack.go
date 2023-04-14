package helper

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{
		slice: []T{},
	}
}

type Stack[T any] struct {
	slice []T
}

func (s *Stack[T]) Push(t T) {
	s.slice = append(s.slice, t)
}

func (s *Stack[T]) Top() T {
	return s.slice[len(s.slice)-1]
}

func (s *Stack[T]) Pop() T {
	top := s.Top()
	s.slice = s.slice[:len(s.slice)-1]
	return top
}

func (s *Stack[T]) ReplaceTop(t T) {
	s.Pop()
	s.Push(t)
}

func (s *Stack[T]) Len() int {
	return len(s.slice)
}

func (s *Stack[T]) Slice() []T {
	return s.slice
}
