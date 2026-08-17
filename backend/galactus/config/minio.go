package config

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIOClient struct {
	Client *minio.Client
	Bucket string
}

func NewMinIOClient(address, user, password, bucket string) (*MinIOClient, error) {
	client, err := minio.New(address, &minio.Options{
		Creds:  credentials.NewStaticV4(user, password, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	exists, errBucketExists := client.BucketExists(ctx, bucket)
	if errBucketExists != nil {
		return nil, errBucketExists
	}
	if !exists {
		err := client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
		if err != nil {
			return nil, err
		}
	}
	return &MinIOClient{
		Client: client,
		Bucket: bucket,
	}, nil
}
