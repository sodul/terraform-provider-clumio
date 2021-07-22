// Copyright 2021. Clumio, Inc.

// Contains mock implementations of the SNS and S3 APIs used for testing.
package clumio

import (
	"bytes"
	"context"
	"io/ioutil"
	"time"

	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sns"

)

// SNSClient is a mock client for SNS to help with unit tests.
type SNSClient struct {

}

// PublishWithContext publishes an SNS notification.
func (s SNSClient) Publish(_ context.Context, _ *sns.PublishInput,
	_ ...func(*sns.Options)) (*sns.PublishOutput, error){
	return &sns.PublishOutput{
		MessageId: aws.String("message-id"),
	}, nil
}

// S3Client is a mock client for S3 to help with unit tests.
type S3Client struct {

}

// HeadObjectWithContext fetches the metadata of a an object from a bucket.
func (s S3Client) HeadObject(_ context.Context, input *s3.HeadObjectInput,
	optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
	currTime := time.Now()
	return &s3.HeadObjectOutput{
		LastModified: &currTime,
	}, nil
}

// GetObjectWithContext fetches an object from a bucket.
func (s S3Client) GetObject(_ context.Context, _ *s3.GetObjectInput,
	_ ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	currTime := time.Now()
	status := StatusObject{
		Status: statusSuccess,
		Reason: aws.String(""),
		Data:   nil,
	}
	statusBytes, err := json.Marshal(status)
	if err != nil{
		return nil, err
	}
	return &s3.GetObjectOutput{
		LastModified: &currTime,
		Body: ioutil.NopCloser(bytes.NewReader(statusBytes)),
	}, nil
}
