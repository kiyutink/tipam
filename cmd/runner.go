package cmd

import "github.com/kiyutink/tipam/tipam"

type runnerFlags struct {
	lock bool
}

var runnerF = runnerFlags{}

func newRunner(p tipam.Persistor) *tipam.Runner {
	opts := []tipam.RunnerOption{}

	if runnerF.lock {
		opts = append(opts, tipam.WithLocking(true))
	}

	return tipam.NewRunner(p, opts...)
}
