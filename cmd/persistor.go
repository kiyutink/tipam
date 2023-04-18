package cmd

import (
	"fmt"

	"github.com/kiyutink/tipam/persist"
	"github.com/kiyutink/tipam/tipam"
)

const persistorLocalYAML = "localyaml"

type persistFlags struct {
	persistorType string

	localYAMLFileName string
}

var persistF = persistFlags{}

func newPersistor() (tipam.Persistor, error) {
	switch persistF.persistorType {
	case persistorLocalYAML:
		p := persist.NewLocalYAML(persistF.localYAMLFileName)
		return p, nil
	}
	return nil, fmt.Errorf("unknown persistor type %v", persistF.persistorType)
}
