package cmd

import (
	"fmt"

	"github.com/kiyutink/tipam/persist"
	"github.com/kiyutink/tipam/tipam"
)

const (
	persistorLocalYAML = "localyaml"
	persistorInMemory  = "inmemory"
)

type persistFlags struct {
	persistor string

	localYAMLFileName string
}

var persistF = persistFlags{}

func newPersistor() (tipam.Persistor, error) {
	switch persistF.persistor {
	case persistorLocalYAML:
		p := persist.NewLocalYAML(persistF.localYAMLFileName)
		return p, nil
	case persistorInMemory:
		p := persist.NewInMemory()
		return p, nil
	}
	return nil, fmt.Errorf("unknown persistor type %v", persistF.persistor)
}
