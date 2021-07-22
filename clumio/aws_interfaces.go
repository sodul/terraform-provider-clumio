// Copyright 2021. Clumio, Inc.

// AWS Interfaces used in the provider.
package clumio

import "context"
import "github.com/aws/aws-sdk-go-v2/service/s3"
import "github.com/aws/aws-sdk-go-v2/service/sns"

// Interface for S3 operations used in the provider.
type S3API interface {
	GetObject(ctx context.Context, params *s3.GetObjectInput,
		optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
	HeadObject(ctx context.Context, params *s3.HeadObjectInput,
		optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error)
}

// Interface for SNS operations used in the provider.
type SNSAPI interface {
	Publish(ctx context.Context, params *sns.PublishInput,
		optFns ...func(*sns.Options)) (*sns.PublishOutput, error)
}
