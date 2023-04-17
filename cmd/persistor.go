package cmd

import (
	"fmt"

	"github.com/kiyutink/tipam/core"
	"github.com/kiyutink/tipam/persist"
)

const persistorLocalYAML = "localyaml"

type persistFlags struct {
	persistorType string

	localYAMLFileName string
}

var persistF = persistFlags{}

func newPersistor() (core.Persistor, error) {
	switch persistF.persistorType {
	case persistorLocalYAML:
		p := persist.NewLocalYAML(persistF.localYAMLFileName)
		return p, nil
	}
	return nil, fmt.Errorf("unknown persistor type %v", persistF.persistorType)
}
