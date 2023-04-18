package tipam

type Persistor interface {
	Persist(*State) error
	Read() (*State, error)
	Lock() error
	Unlock() error
}

type Runner struct {
	persistor Persistor
	doLock    bool
}

type runnerParams struct {
	// doLock specifies whether to use persistor's locks
	// when operating on the state
	doLock bool
}

type RunnerOption func(*runnerParams)

func WithLocking(do bool) RunnerOption {
	return func(rp *runnerParams) {
		rp.doLock = do
	}
}

func NewRunner(p Persistor, opts ...RunnerOption) *Runner {
	r := &Runner{
		persistor: p,
	}

	params := &runnerParams{
		doLock: true,
	}

	for _, opt := range opts {
		opt(params)
	}

	r.doLock = params.doLock

	return r
}

func (r *Runner) ReadState() (*State, error) {
	return r.persistor.Read()
}

func (r *Runner) PersistState(state *State) error {
	return r.persistor.Persist(state)
}
