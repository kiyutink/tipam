package tipam

// Persistor is used to Write/Read/Lock/Unlock state to a storage
type Persistor interface {
	Persist(*State) error
	Read() (*State, error)
	Lock() error
	Unlock() error
}

// Runner is the main structure that defines all the externally-available methods.
// In order to use tipam, one must initialize an instance of runner.
// Do not initialize Runner directly, use the constructor function tipam.NewRunner instead
type Runner struct {
	persistor Persistor
	doLock    bool
}

type RunnerOpts struct {
	// DoLock specifies whether to use persistor's locks
	// when operating on the state
	DoLock bool
}

// Initializes a new Runner instance and returns a pointer to it. Only use NewRunner
// for initializing a Runner, do not initialize it directly
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
