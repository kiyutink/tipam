package cmd

import (
	"fmt"
	"time"

	"github.com/kiyutink/tipam/persist"
	"github.com/kiyutink/tipam/tipam"
)

const (
	persistorLocalYAML = "localyaml"
	persistorInMemory  = "inmemory"
	persistorS3Dynamo  = "s3dynamo"
)

type persistFlags struct {
	persistor string

	localYAMLFileName string

	s3DynamoBucket        string
	s3DynamoKeyInBucket   string
	s3DynamoTable         string
	s3DynamoLeaseDuration int
	s3DynamoPollInterval  int
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
	case persistorS3Dynamo:
		p, err := persist.NewS3Dynamo(
			persistF.s3DynamoBucket,
			persistF.s3DynamoKeyInBucket,
			persistF.s3DynamoTable,
			time.Duration(persistF.s3DynamoLeaseDuration)*time.Second,
			time.Duration(persistF.s3DynamoPollInterval)*time.Second,
		)
		if err != nil {
			return nil, err
		}
		return p, nil
	}
	return nil, fmt.Errorf("unknown persistor type %v", persistF.persistor)
}
