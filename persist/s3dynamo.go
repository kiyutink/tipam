package persist

import (
	"bytes"
	"context"
	"errors"
	"time"

	"cirello.io/dynamolock/v2"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"gopkg.in/yaml.v3"

	"github.com/kiyutink/tipam/tipam"
)

type S3Dynamo struct {
	s3client   *s3.Client
	lock       *dynamolock.Lock
	lockClient *dynamolock.Client

	bucket        string
	keyInBucket   string
	table         string
	leaseDuration time.Duration
	pollInterval  time.Duration
}

func NewS3Dynamo(bucket string, keyInBucket string, table string, leaseDuration time.Duration, pollInterval time.Duration) (*S3Dynamo, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}
	s3Client := s3.NewFromConfig(cfg)

	return &S3Dynamo{
		s3client: s3Client,

		bucket:        bucket,
		keyInBucket:   keyInBucket,
		table:         table,
		leaseDuration: leaseDuration,
		pollInterval:  pollInterval,
	}, nil
}

func (s3d *S3Dynamo) Persist(s *tipam.State) error {
	yamlState := newEmptyYAMLState()
	for cidr, r := range s.Claims {
		yamlState.Claims[cidr] = yamlStateClaim{Tags: r.Tags, Final: r.Final}
	}

	buf := bytes.NewBuffer([]byte{})
	encoder := yaml.NewEncoder(buf)
	encoder.SetIndent(2)
	err := encoder.Encode(&yamlState)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	_, err = s3d.s3client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s3d.bucket),
		Key:    aws.String(s3d.keyInBucket),
		Body:   buf,
	})

	if err != nil {
		return err
	}

	return nil
}

func (s3d *S3Dynamo) Read() (*tipam.State, error) {
	ys := newEmptyYAMLState()

	s3Obj, err := s3d.s3client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(s3d.bucket),
		Key:    aws.String(s3d.keyInBucket),
	})

	var noSuchKey *types.NoSuchKey
	if errors.As(err, &noSuchKey) {
		return tipam.NewState(), nil
	}

	err = decodeYAMLState(ys, s3Obj.Body)
	if err != nil {
		return nil, err
	}

	s, err := yamlStateToState(ys)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s3d *S3Dynamo) Lock() error {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return err
	}

	c, err := dynamolock.New(
		dynamodb.NewFromConfig(cfg),
		s3d.table,
		dynamolock.WithLeaseDuration(s3d.leaseDuration),
		dynamolock.WithHeartbeatPeriod(s3d.pollInterval),
	)
	if err != nil {
		panic(err)
	}
	s3d.lockClient = c

	l, err := c.AcquireLock(
		"lock",
		dynamolock.WithRefreshPeriod(s3d.pollInterval),
		dynamolock.WithAdditionalTimeToWaitForLock(s3d.leaseDuration),
	)
	s3d.lock = l
	if err != nil {
		panic(err)
	}

	return nil
}

func (s3d *S3Dynamo) Unlock() error {
	_, err := s3d.lockClient.ReleaseLock(s3d.lock)
	if err != nil {
		panic(err)
	}
	return nil
}
