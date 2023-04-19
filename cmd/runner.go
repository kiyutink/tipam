package cmd

import "github.com/kiyutink/tipam/tipam"

type runnerFlags struct {
	lock bool
}

var runnerF = runnerFlags{}

// newRunner returns a new *tipam.Runner with runnerFlags passed in as options
func newRunner(p tipam.Persistor) *tipam.Runner {
	return tipam.NewRunner(p, tipam.RunnerOpts{DoLock: runnerF.lock})
}
