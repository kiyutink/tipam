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

type RunnerOpts struct {
	// DoLock specifies whether to use persistor's locks
	// when operating on the state
	DoLock bool
}

func NewRunner(p Persistor, opts RunnerOpts) *Runner {
	r := &Runner{
		persistor: p,
	}

	r.doLock = opts.DoLock

	return r
}

func (r *Runner) ReadState() (*State, error) {
	return r.persistor.Read()
}

func (r *Runner) PersistState(state *State) error {
	return r.persistor.Persist(state)
}
