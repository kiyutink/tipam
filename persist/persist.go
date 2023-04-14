package persist

import (
	"fmt"

	"github.com/kiyutink/tipam/core"
)

type PersistFlags struct {
	PersistorType string

	LocalYAMLFileName string
}

func NewPersistor(pf PersistFlags) (core.Persistor, error) {
	switch pf.PersistorType {
	case PersistorLocalYaml:
		p := NewLocalYAMLPersistor(pf.LocalYAMLFileName)
		return p, nil
	}
	return nil, fmt.Errorf("unknown persistor type %v", pf.PersistorType)
}
