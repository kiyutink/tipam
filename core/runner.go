package core

type Persistor interface {
	Persist(*State) error
	Read() (*State, error)
}

type Runner struct {
	Persistor Persistor
}

func NewRunner(p Persistor) *Runner {
	return &Runner{
		Persistor: p,
	}
}
